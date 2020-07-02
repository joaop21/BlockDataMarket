const express = require('express');
const router = express.Router();
const app = require('../app')
const database = require('../scripts/database');
const wiki = require('../scripts/wikiUtils')
const multer = require('multer')
const path = require('path');
const storage = multer.diskStorage({
    destination: function (req, file, cb) {
        cb(null, path.join(__dirname, '../uploads'))
    },
    filename: function (req, file, cb) {
        cb(null, file.originalname)
    }
})
const upload = multer({ storage: storage })
const crypto = require('../scripts/crypto')

var chaincode;

router.use(async function (req, res, next) {
    chaincode = await app.getChaincode()
    if (chaincode != null)
        next()
    else res.send({ result: "!ok" });
})

/* GET announcements */
router.get('/', async function (req, res) {
    let category = req.query.category
    let ownerId = req.query.ownerId
    let lt = req.query.lt

    try {
        if (category && lt)
            result = await chaincode.evaluateTransaction('AnnouncementContract:GetAnnouncementsByCategoryLowerThan', category, lt)
        else if (category)
            result = await chaincode.evaluateTransaction('AnnouncementContract:GetAnnouncementsByCategory', category)
        else if (ownerId)
            result = await chaincode.evaluateTransaction('AnnouncementContract:GetAnnouncementsByOwner', ownerId)
        else if (lt)
            result = await chaincode.evaluateTransaction('AnnouncementContract:GetAnnouncementsLowerThan', lt)
        else
            result = await chaincode.evaluateTransaction('AnnouncementContract:GetAnnouncements')
    }
    catch (err) {
        res.send({ error: err.toString() })
    }

    res.send({ result: JSON.parse(result) });
});

/* POST announcement */
router.post('/', upload.single('data_file'), async function (req, res) {
    let file = req.file
    let queries = req.body.queries
    let category = req.body.category

    if (file && queries && category) {
        try{
            queriesArray = JSON.parse(queries)

            const dataId = await database.putContent(file.path);
            const queryPrices = await wiki.getQueryPrices(dataId, queriesArray);
            let prices = JSON.parse(queryPrices)

            var announcement = await chaincode.submitTransaction('AnnouncementContract:MakeAnnouncement', dataId, queries, queryPrices, category)
            announcement = JSON.parse(announcement)

            res.send({ result: announcement });
            if (announcement != null) {
                const eventName1 = 'Query:' + announcement.announcementId;
                const eventName2 = 'Update:' + announcement.announcementId;
                const listener = async (event) => {
                    // query arrives
                    if (event.eventName === eventName1) {
                        event = event.payload.toString();
                        event = JSON.parse(event);
                        // putResponseLogic
                        const queryIndex = queriesArray.indexOf(event.query);
                        const response = await wiki.getResponseContent(dataId, event.query, prices[queryIndex], event.price);
                        const issuer = await chaincode.evaluateTransaction('IdentificationContract:GetIdentification', event.issuerId);
                        const issuerJson = JSON.parse(issuer);
                        const criptogram = crypto.encrypt(response, issuerJson.publicKey);
                        return await chaincode.submitTransaction('QueryContract:PutResponse', event.queryId, criptogram);
                    }
                    // changes arrives
                    else if (event.eventName === eventName2) {
                        event = event.payload.toString();
                        event = JSON.parse(event);
                        // change pricesArray
                        prices = event.prices
                    }
                };
                await chaincode.addContractListener(listener);
            }
        } catch(err) {
            res.status(400).send({ error: err.toString() });
        }
    }
    else 
        res.status(400).send({ error: "You must provide a file, its category and prices." })
});

/* POST announcement */
router.post('/UpdatePrices', upload.none(), async function (req, res) {
    let announcementId = req.body.announcementId
    let updates = req.body.updates

    try {
        if (announcementId){
            result = await chaincode.submitTransaction('AnnouncementContract:UpdateQueryPrices', announcementId, updates)
            res.send({ result: JSON.parse(result) });
        }
    }
    catch (err) {
        res.send({ error: err.toString() })
    }

});

module.exports = router;