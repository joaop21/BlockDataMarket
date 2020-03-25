# Deploy and Start Chaincode into the Network

There are some assumptions that we made when running this script:
  - The fabric-samples folder exists and is in BlockDataMarket/fabric-samples/ ;
  - The test-network folder is BlockDataMarket/fabric-samples/test-network ;
  - The test-network must be up-and-running;
  - The chaincode folder is BlockDataMarket/chaincode/ ;

You should run this command in order to run the script and also to export the file variables into session variables:
```
source ./deploy-chaincode.sh [Flags]
```

This script has a moment of interaction where it needs an input, be aware.

ATTENTION: this script also initializes the chaincode with 'AnnouncementContract:Instantiate' function.
