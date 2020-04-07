const assert = require('assert');
const mongo = require('mongodb');
const { v1: uuidv1 } = require('uuid');
const config = require('./config.json');
const fs = require('fs');

// Connection URL
const MongoClient = mongo.MongoClient(config.database.url, { useUnifiedTopology: true });

// only in this module scope
let callback, db;

// Use connect method to connect to Mongo
// This method connects to DB asynchronously
MongoClient.connect(function(err, client) {
  assert.equal(null, err);
  db = client.db(config.database.name);
  callback(db);
});

// gives a DB client when connected
var connect = new Promise(function (returnedValue){
  if(typeof db != 'undefined'){
    returnedValue(db)
  } else {
    // returnedValue is an anonymous function, so it takes the value of callback when it's available
    callback = returnedValue;
  }
});


// exports this function to put data in DB
module.exports.putContent = function(filepath){
  return new Promise(async (dataID) => {
    let dbconn = await connect;
    const content = fs.readFileSync(filepath, 'utf8');
    const uuid1 = uuidv1();
    const data = {
      dataID: uuid1,
      filepath: filepath,
      content: content
    };
    dbconn.collection(config.database.collection).insertOne(data);
    dataID(uuid1);
  });
};

// exports this function to get data from DB
module.exports.getContent = function(dataID){
  return new Promise(async (res) => {
    let dbconn = await connect;
    var criteria = {
      dataID: dataID
    };
    dbconn.collection(config.database.collection).findOne(criteria, function(err, result) {
      if (err) throw err;
      res(result);
    });
  });
};
