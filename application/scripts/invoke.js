'use strict';

const { Gateway, Wallets, ContractListener } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const database = require('./database');

let contract;

async function makeAnnouncement(funcName, filename, prices, category){
    const dataId = await database.putContent(filename);
    console.log(dataId + " " + prices + " " + category);
    const announcementId = await contract.submitTransaction(funcName, dataId, prices, category);
    if(announcementId != null){
        const eventName = 'Query:' + announcementId;
            const listener = async (event) => {
                console.log('Received event: '+event);
                if (event.eventName === eventName) {
                    console.log('Event Payload: '+ event.payload.toString('utf8'))
                }
            };
        await contract.addContractListener(listener);
    }
    return announcementId;
}

//Prototype to check query sintax
function checkQuerySintax(query){
    if(1==0)
        return [false, "Error: Invalid query syntax"];
    
    return [true, null]
}

async function makeQuery(funcName, announcementId, queryArg, price){
    const announcement = await contract.submitTransaction('AnnouncementContract:GetAnnouncement', announcementId);
    if(announcement){
        //check querySyntax
        const check = checkQuerySintax(queryArg);
	    if(check[0]){
            return (await contract.submitTransaction(funcName, announcementId, queryArg, price));
        }else{
            return check[1];
        }
    }else{
        return "Error: Invalid Announcement ID";
    }
}

//ir buscar a resposta ao ficheiro na bd, truncar conforme o nivel
async function getResponse(dataId, level){
    const content = await database.getContent(dataId);
    const filePercentage = 0.5*level; //level 1 a 2
    return content.slice(0, content.length*filePercentage);
}


async function putResponse(funcName, queryid){
    const query = await contract.submitTransaction('QueryContract:GetQuery',queryid);
    if(query){
        const queryJson = JSON.parse(query);
        const announcementId = queryJson.announcementId;
        
        const announcement = await contract.submitTransaction('AnnouncementContract:GetAnnouncement', announcementId);
        const announcementJson = JSON.parse(announcement);
        
        const prices = announcementJson.prices; 
        const index = prices.findIndex( (price) => price === queryJson.price);
        const response = index !== -1
            ? await getResponse(announcementJson.dataId, index + 1)
            : "Offer declined, price didn't match any of the levels";
        return await contract.submitTransaction(funcName, queryid, response);
    }else{
        return "Error: Query doesn't exist"
    }

}

async function main() {
    try {
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
        contract = network.getContract('dataMarket');

        // accept args from stdin
        const args = process.argv.slice(2);
        let result = null;
        // submit transaction depending on first arg
        switch (args[0]) {
            case 'AnnouncementContract:MakeAnnouncement':
                result = await makeAnnouncement(args[0], args[1], args[2], args[3]);
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
                result = await makeQuery(args[0], args[1], args[2], args[3]);
                break;
            case 'QueryContract:PutResponse':
                result = await putResponse(args[0], args[1]);
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
        // await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main().then(()=>{
    console.log('done');
}).catch((e)=>{
    console.log('Final error checking.......');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);
});
