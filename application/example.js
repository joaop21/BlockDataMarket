const mongo = require('./mongo_connection');

// example of connection
mongo.connect((db) => {
  db.collection('Category').find().toArray((err,results) => {
    console.log(results);
  });
})
