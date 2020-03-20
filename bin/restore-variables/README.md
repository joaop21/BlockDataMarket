# Restore the Environment Network Variables

This script is useful because when you exit the machine, you lost all the network variables that were set. With this script some variables are redefined and after you run it, you're able to invoke functions on Org1.

There are some assumptions that we made when running this script:
  - The fabric-samples folder exists and is in BlockDataMarket/fabric-samples/ ;
  - The test-network folder is BlockDataMarket/fabric-samples/test-network ;
  - The test-network must be up-and-running;

You should run this command in order to run the script and also to export the file variables into session variables:
```
source ./restore-variables.sh
```
