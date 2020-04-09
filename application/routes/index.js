var express = require('express');
var router = express.Router();
var package = require('../package.json');

/* GET home page. */
router.get('/', function(req, res, next) {
  res.send({ app: package.name, description: package.description });
});

module.exports = router;
