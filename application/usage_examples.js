const mongo = require('./database');

// for putting content in db
mongo.putContent('./package.json').then((dataID) => {
  console.log(dataID)
});

// for putting content in db
mongo.getContent('d7e61b40-78ef-11ea-bb13-f169cd2409f9').then((result) => {
  console.log(result)
});
