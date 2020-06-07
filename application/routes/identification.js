const express = require('express');
const router = express.Router();
const app = require('../app');
const database = require('../scripts/database');
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

/* GET identification */
router.get('/', async function (req, res) {
    var identificationId = req.query.identificationId

    var result;
    if (identificationId){
        result = await chaincode.submitTransaction('IdentificationContract:GetIdentification', identificationId);
        res.send({ result: JSON.parse(result) });
    }
    else {
        res.status(400).send({error : "No identification Id was provided"})
    }

});

/* POST identification */
router.post('/', upload.none(), async function (req, res) {
    var name =  req.body.name

    if (name) {
        try{
            var publicKey = crypto.generateKeys();

            var identification = await chaincode.submitTransaction('IdentificationContract:MakeIdentification', name, publicKey);
        
            res.send({ result: JSON.parse(identification) });
        

        } catch(err) {
            res.status(400).send({ error: err.toString() });
        }
    }
    else 
        res.status(400).send({ error: "You must provide name in order to make an identification" })
});


//Prototype to check query sintax
function checkQuerySintax(query) {
    if (1 == 0)
        return false;

    return true;
}

module.exports = router;