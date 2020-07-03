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

/* GET category */
router.get('/', async function (req, res) {
    var categoryName = req.query.categoryName

    var result;
    try {
        if (categoryName) {
            result = await chaincode.evaluateTransaction('CategoryContract:GetCategory', categoryName);
        }
        else {
            result = await chaincode.evaluateTransaction('CategoryContract:GetCategories');
        }
        res.send({ result: JSON.parse(result) });
    }
    catch (err) {
        res.send({ error: err.toString() })
    }
});

/* POST category */
router.post('/', upload.none(), async function (req, res) {
    var name = req.body.name
    var queries = req.body.queries

    if (name && queries) {
        try {
            var result = await chaincode.submitTransaction('CategoryContract:MakeCategory', name, queries);
            res.send({ result: result.toString() });
        } catch (err) {
            res.status(400).send({ error: err.toString() });
        }
    }
    else
        res.status(400).send({ error: "You must provide a name and a list of possible queries in order to make a category" })
});


//Prototype to check query sintax
function checkQuerySintax(query) {
    if (1 == 0)
        return false;

    return true;
}

module.exports = router;