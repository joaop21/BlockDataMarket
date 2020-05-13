const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function getContract() {
    // load the network configuration
    const ccpPath = path.resolve(__dirname, '..', '..', "fabric-samples", "test-network", "organizations",
        "peerOrganizations", "org1.example.com", 'connection-org1.json');
    let ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const identity = await wallet.get('admin');
    if (!identity) {
        console.log('An identity for the user "admin" does not exist in the wallet');
        console.log('Run the enrollAdmin.js application before retrying');
        return;
    }

    // Create a new gateway for connecting to our peer node.
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: true, asLocalhost: true } });

    // Get the network (channel) our contract is deployed to.
    const network = await gateway.getNetwork('mychannel');

    // Get the chaincode from the network.
    return network.getContract('dataMarket');
}

module.exports = getContract;