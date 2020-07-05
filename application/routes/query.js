const express = require('express');
const router = express.Router();
const app = require('../app');
const multer = require('multer')
const upload = multer()
const crypto = require('../scripts/crypto')
var chaincode;

router.use(async function (req, res, next) {
    chaincode = await app.getChaincode()
    if (chaincode != null)
        next()
    else res.send({ result: "!ok" });
})

/* GET query */
router.get('/', async function (req, res) {
    var announcementId = req.query.announcementId
    var queryId = req.query.queryId
    var issuerId = req.query.issuerId

    var result;
    try {
        if (queryId) {
            result = await chaincode.evaluateTransaction('QueryContract:GetQuery', queryId);

            const resultJson = JSON.parse(result);

            let announcement
            try {
                announcement = await chaincode.evaluateTransaction('AnnouncementContract:GetAnnouncement', resultJson.announcementId);
            } catch (err) {
                res.status(400).send({ Error: "Invalid Announcement ID" });
            }

            const cryptogram = resultJson.response;
            const announcementJson = JSON.parse(announcement);
            const owner = await chaincode.evaluateTransaction('IdentificationContract:GetIdentification', announcementJson.ownerId);
            const ownerJson = JSON.parse(owner);
            resultJson.response = crypto.decrypt(cryptogram, ownerJson.publicKey);
            result = JSON.stringify(resultJson)

        }
        else if (announcementId)
            result = await chaincode.evaluateTransaction('QueryContract:GetQueriesByAnnouncement', announcementId);
        else if (issuerId)
            result = await chaincode.evaluateTransaction('QueryContract:GetQueriesByIssuer', issuerId);
        else
            res.status(400).send({ error: "Neither query, announcement or issuer Id was provided. You must provide one of them as an argument." })

        res.send({ result: JSON.parse(result) });
    }
    catch (err) {
        res.status(400).send({ error: err.toString() });
    }
});

/* POST query */
router.post('/', upload.none(), async function (req, res) {
    var announcementId = req.body.announcementId
    var query = req.body.query
    var price = req.body.price

    var announcement;
    try {
        announcement = await chaincode.evaluateTransaction('AnnouncementContract:GetAnnouncement', announcementId);
    } catch (err) {
        res.status(400).send({ Error: "Invalid Announcement ID" });
    }

    if (announcement) {
        var query;
        try{
            query = await chaincode.submitTransaction("QueryContract:MakeQuery", announcementId, query, price);
            query = JSON.parse(query.toString())
        }
        catch (err){
            res.status(400).send({ Error: err.toString() });
        }

        if (query != null) {
            console.log( 'Response:' + query.queryId)
            const eventName = 'Response:' + query.queryId;
            const listener = async (event) => {
                if (event.eventName === eventName) {
                    event = event.payload.toString();
                    event = JSON.parse(event);
                    const cryptogram = event.response;
                    const announcementJson = JSON.parse(announcement);
                    const owner = await chaincode.evaluateTransaction('IdentificationContract:GetIdentification', announcementJson.ownerId);
                    const ownerJson = JSON.parse(owner);
                    event.response = crypto.decrypt(cryptogram, ownerJson.publicKey);
                    res.send({ result: event });
                }
            };
            await chaincode.addContractListener(listener);
        }
    }
});



module.exports = router;