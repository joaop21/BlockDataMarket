'use strict';

const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const mycrypto = require('./crypto')
const database = require('./database');
const wiki = require('./wikiUtils')

let contract;


async function makeAnnouncement(filename, category, queries) {
    const dataId = await database.putContent(filename);

    const prices = wiki.getQueryPrices(filename);
    const pricesArray = JSON.parse(prices)

    const announcementId = await contract.submitTransaction('AnnouncementContract:MakeAnnouncement', dataId, queries, prices, category);

    if (announcementId != null) {
        const eventName = 'Query:' + announcementId;
        const listener = async (event) => {
            if (event.eventName === eventName) {
                event = event.payload.toString();
                event = JSON.parse(event);
                // putResponseLogic
                const queryIndex = queries.indexOf(event.query);
                const response = wiki.getResponseContent(filename, event.query, pricesArray[queryIndex], event.price);
                const issuer = await contract.submitTransaction('IdentificationContract:GetIdentification', event.issuerId);
                const issuerJson = JSON.parse(issuer);
                const criptogram = mycrypto.encrypt(response, issuerJson.publicKey);
                return await contract.submitTransaction('QueryContract:PutResponse', event.queryId, criptogram);
            }
        };
        await contract.addContractListener(listener);
    }
    return announcementId;
}

//Prototype to check query sintax
function checkQuerySintax(query) {
    if (1 == 0)
        return [false, "Error: Invalid query syntax"];

    return [true, null]
}

async function makeQuery(funcName, announcementId, queryArg, price) {
    const announcement = await contract.submitTransaction('AnnouncementContract:GetAnnouncement', announcementId);
    if (announcement) {
        //check querySyntax
        const check = checkQuerySintax(queryArg);
        if (check[0]) {
            const queryId = await contract.submitTransaction(funcName, announcementId, queryArg, price);
            if (queryId != null) {
                const eventName = 'Response:' + queryId;
                const listener = async (event) => {
                    if (event.eventName === eventName) {
                        event = event.payload.toString();
                        event = JSON.parse(event);
                        const cryptogram = event.response;
                        const announcementJson = JSON.parse(announcement);
                        const owner = await contract.submitTransaction('IdentificationContract:GetIdentification', announcementJson.ownerId);
                        const ownerJson = JSON.parse(owner);
                        const plaintext = mycrypto.decrypt(cryptogram, ownerJson.publicKey);
                        console.log('Received Response: ' + plaintext);
                    }
                };
                await contract.addContractListener(listener);
            }
            return queryId;
        } else {
            return check[1];
        }
    } else {
        return "Error: Invalid Announcement ID";
    }
}

//ir buscar a resposta ao ficheiro na bd, truncar conforme o nivel
async function getResponse(dataId, level) {
    const content = await database.getContent(dataId);
    const filePercentage = 0.5 * level; //level 1 a 2
    return content.slice(0, content.length * filePercentage);
}

//deprecated after events implementation
async function putResponse(funcName, queryid) {
    const query = await contract.submitTransaction('QueryContract:GetQuery', queryid);
    if (query) {
        const queryJson = JSON.parse(query);
        const announcementId = queryJson.announcementId;

        const announcement = await contract.submitTransaction('AnnouncementContract:GetAnnouncement', announcementId);
        const announcementJson = JSON.parse(announcement);

        const prices = announcementJson.prices;
        const index = prices.findIndex((price) => price === queryJson.price);
        const response = index !== -1
            ? await getResponse(announcementJson.dataId, index + 1)
            : "Offer declined, price didn't match any of the levels";
        //encrypt response
        const issuer = await contract.submitTrabsaction('IdentificationContract:GetIdentification', queryJson.issuerId);
        const issuerJson = JSON.parse(issuer)
        let criptogram = mycrypto.encrypt(response, issuerJson.publicKey)
        return await contract.submitTransaction(funcName, queryid, criptogram);
    } else {
        return "Error: Query doesn't exist"
    }

}


async function makeIdentification(funcName, name, ip) {
    let publicKey = mycrypto.generateKeys();

    await contract.submitTransaction(funcName, name, ip, publicKey);

    return "Your Private-Key is saved under priv.pem, keep it save"
}


//deprecated after events implementation
async function getQuery(funcName, queryId) {
    query = await contract.submitTransaction(funcName, queryId);
    const queryJson = JSON.parse(query);
    let announcement = await contract.submitTransaction('AnnouncementContract:GetAnnouncement', queryJson.announcementId);
    const announcementJson = JSON.parse(announcement);
    let owner = await contract.submitTransaction('IdentificationContract:GetIdentification', announcementJson.ownerId);
    let ownerJson = JSON.parse(owner)

    return mycrypto.decrypt(responseCriptogram, owner.publicKey)
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
        contract = network.getContract('dataMarket', 1);

        // accept args from stdin
        const args = process.argv.slice(2);
        let result = null;
        // submit transaction depending on first arg*/
        switch (args[0]) {
            case 'AnnouncementContract:MakeAnnouncement':
                result = await makeAnnouncement(args[1], args[2], args[3]);
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
            case 'QueryContract:GetQueriesByAnnouncement':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'QueryContract:GetQueriesByIssuer':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'IdentificationContract:MakeIdentification':
                result = await makeIdentification(args[0], args[1], args[2]);
                break;
            case 'IdentificationContract:GetIdentification':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'CategoryContract:MakeCategory':
                result = await contract.submitTransaction(args[0], args[1], args[2]);
                break;
            case 'CategoryContract:GetCategory':
                result = await contract.submitTransaction(args[0], args[1]);
                break;
            case 'CategoryContract:GetCategories':
                result = await contract.submitTransaction(args[0]);
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

main().then(() => {
    console.log('done');
}).catch((e) => {
    console.log('Final error checking.......');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);
});
