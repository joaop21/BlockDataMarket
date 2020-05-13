var express = require('express');
var createError = require('http-errors');
var cookieParser = require('cookie-parser');
var logger = require('morgan');
var getContract = require('./scripts/contract')

var indexRouter = require('./routes/index');
var announcementRouter = require('./routes/announcement');
var queryRouter = require('./routes/query');
var identificationRouter = require('./routes/identification');


var chaincode = null;
module.exports.getChaincode = async function getChaincode() {
  if (chaincode == null)
    chaincode = await getContract();
  return chaincode;  
}

var app = express();
app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());

app.disable('etag'); // serve para impedir que sejam enviadas respostas com HTTP Status 304

app.use('/', indexRouter);
app.use('/announcements', announcementRouter);
app.use('/query', queryRouter);
app.use('/identification', identificationRouter);



// catch 404 and forward to error handler
app.use(function(req, res, next) {
  next(createError(404));
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.send(err.message);
});

module.exports = {
  app: app,
};
