const assert = require('assert');
const mongo = require('mongodb');
const { v1: uuidv1 } = require('uuid');
const config = require('../config.json');
const fs = require('fs');

// Connection URL
const MongoClient = mongo.MongoClient(config.database.url, { useUnifiedTopology: true });

// only in this module scope
let callback, db;

// Use connect method to connect to Mongo
// This method connects to DB asynchronously

MongoClient.connect().then(client => {
  db = client.db(config.database.name);
  callback(db);
}).catch(err => {
  console.error(err)
});

// gives a DB client when connected
const connect = new Promise((returnedValue) => {
  if (typeof db != 'undefined') {
    returnedValue(db)
  } else {
    // returnedValue is an anonymous function, so it takes the value of callback when it's available
    callback = returnedValue;
  }
});


// exports this function to put data in DB
module.exports.putContent = async function (filepath) {
  let dbconn = await connect;
  const content = fs.readFileSync(filepath, 'utf8');
  const uuid1 = uuidv1();
  const data = {
    dataID: uuid1,
    filepath: filepath,
    content: content
  };
  dbconn.collection(config.database.collection).insertOne(data);
  return uuid1
};

// exports this function to get data from DB
module.exports.getContent = (dataID) => {
  return new Promise(async (res) => {
    let dbconn = await connect;
    const criteria = {
      dataID: dataID
    };

    dbconn.collection(config.database.collection).findOne(criteria)
        .then(result => res(result.content))
        .catch(err => console.error(err));
  });
};
