const express = require('express');
const router = express.Router();
const app = require('../app')
const database = require('../scripts/database');
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

    console.log(category)
    console.log(ownerId)
    console.log(lt)

    try {
        if (category)
            result = await chaincode.submitTransaction('AnnouncementContract:GetAnnouncementsByCategory', category)
        else if (ownerId)
            result = await chaincode.submitTransaction('AnnouncementContract:GetAnnouncementsByOwner', ownerId)
        else if (lt)
            result = await chaincode.submitTransaction('AnnouncementContract:GetAnnouncementsLowerThan', lt)
        else
            result = await chaincode.submitTransaction('AnnouncementContract:GetAnnouncements')
    }
    catch (err) {
        res.send({ error: err.toString() })
    }

    res.send({ result: result.toString() });
});

/* POST announcement */
router.post('/', upload.single('data_file'), async function (req, res) {
    let file = req.file
    let prices = JSON.parse(req.body.prices)
    let category = req.body.category

    console.log(req.body)
    console.log(req.file)

    const dataId = await database.putContent(file.path);

    console.log(dataId + " " + prices + " " + category);

    result = await chaincode.submitTransaction('AnnouncementContract:MakeAnnouncement', dataId, prices, category)

    res.send({ result: result.toString() });

    if (file && prices && category) {
        try{
            const dataId = await database.putContent(file.path);
            console.log(dataId + " " + prices + " " + category);
            result = await chaincode.submitTransaction('AnnouncementContract:MakeAnnouncement', dataId, prices, category)
            res.send({ result: result.toString() });
        } catch(err) {
            res.send({ error: err.toString() });
        }
    }
    else 
        res.status(400).send({ error: "You must provide a file, its category and prices." })
});

module.exports = router;