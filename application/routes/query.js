const express = require('express');
const router = express.Router();
const app = require('../app');
const multer = require('multer')
const upload = multer()

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

    console.log(announcementId)
    console.log(queryId)
    console.log(issuerId)

    var result;
    try {
        if (queryId)
            result = await chaincode.submitTransaction('QueryContract:GetQuery', queryId);
        else if (announcementId)
            result = await chaincode.submitTransaction('QueryContract:GetQueriesByAnnouncement', announcementId);
        else if (issuerId)
            result = await chaincode.submitTransaction('QueryContract:GetQueriesByIssuer', issuerId);
        else
            res.status(400).send({ error: "Neither query, announcement or issuer Id was provided. You must provide one of them as an argument." })

        res.send({ result: result.toString() });
    }
    catch (err) {
        res.send({ error: err.toString() });
    }
});

/* POST query */
router.post('/', upload.none(), async function (req, res) {
    var announcementId = req.body.announcementId
    var queryArg = req.body.queryArg
    var price = req.body.price

    console.log(req.body)

    var announcement;
    try {
        announcement = await chaincode.submitTransaction('AnnouncementContract:GetAnnouncement', announcementId);
    } catch (err) {
        res.send({ Error: "Invalid Announcement ID" });
    }

    if (announcement) {
        const check = checkQuerySintax(queryArg);
        if (check) {
            var result = await chaincode.submitTransaction("QueryContract:MakeQuery", announcementId, queryArg, price);
            res.send({ result: result.toString() });
        } else {
            res.send({ Error: "Invalid Query Syntax" });
        }
    }
});


//Prototype to check query sintax
function checkQuerySintax(query) {
    if (1 == 0)
        return false;

    return true;
}

module.exports = router;