const assert = require('assert');
const mongo = require('mongodb')

// Connection URL
const url = 'mongodb://localhost:27017';
const MongoClient = mongo.MongoClient(url, { useUnifiedTopology: true });

// Database Name
const dbName = 'Org1DB';

// only in this module scope
let callback, db;

// Use connect method to connect to Mongo
MongoClient.connect(function(err, client) {
  assert.equal(null, err);
  console.log("Connected successfully to Mongo");

  db = client.db(dbName);
  callback(db);
});

// exports db when connected
module.exports.connect = (returnedValue) => {
  if(typeof db != 'undefined'){
    returnedValue(db)
  } else {
    callback = returnedValue;
  }
}
