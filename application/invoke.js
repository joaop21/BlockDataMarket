'use strict';

const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function main() {
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', "fabric-samples", "test-network", "organizations",
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
        const contract = network.getContract('dataMarket');

        // accept args from stdin
        var args = process.argv.slice(2);
        let result = null;
        // submit transaction depending on first arg
        switch (args[0]) {
            case 'AnnouncementContract:MakeAnnouncement':
                result = await contract.submitTransaction(args[0], args[1], args[2], args[3], args[4]);
                break;
            case 'AnnouncementContract:GetAnnouncements':
                result = await contract.submitTransaction(args[0]);
                break;
            case 'AnnouncementContract:GetAnnouncementsByCategory':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'AnnouncementContract:GetAnnouncementsByOwner':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'AnnouncementContract:GetAnnouncementsLowerThan':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'QueryContract:MakeQuery':
                result = await contract.submitTransaction(args[0], args[1], args[2], args[3], args[4]);
                break;
            case 'QueryContract:PutResponse':
                result = await contract.submitTransaction(args[0], args[1], args[2], args[3], args[4]);
                break;
            case 'QueryContract:GetQueriesByAnnouncement':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'QueryContract:GetQueriesByIssuer':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'IdentificationContract:MakeIdentification':
                result = await contract.submitTransaction(args[0], args[1], args[2], args[3]);
                break;
            case 'IdentificationContract:GetIdentification':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
        }

        console.log(`Transaction has been submitted, result is: ${result.toString()}`);

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
