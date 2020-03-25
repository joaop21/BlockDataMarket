#!/bin/bash

set -e
. "../functions.sh"

# Print the usage message
function printHelp() {
  echo "Usage: "
  echo "  deploy-chaincode.sh [Flags]"
  echo "    Flags:"
  echo "    -pn <package name> - defines a name for the smart contract package"
  echo "    -cn <channel name> - sets the channel name if it isn't 'mychannel'"
  echo "    -h (print this message)"
  echo
  echo "Taking all defaults:"
  echo "  deploy-chaincode.sh"
}

# default values
PACKAGE_NAME=dataMarket
CHANNEL_NAME=mychannel

# parse flags

while [[ $# -ge 1 ]] ; do
  key="$1"
  case $key in
  -h )
    printHelp
    exit 0
    ;;
  -pn )
    PACKAGE_NAME="$2"
    shift
    ;;
  -cn )
      CHANNEL_NAME="$2"
      shift
      ;;
  * )
    echo
    echo "Unknown flag: $key"
    echo
    printHelp
    exit 1
    ;;
  esac
  shift
done

pp_info "Deploy" "Package the smart contract"
cd ../../chaincode/
rm go.mod go.sum
go mod init $PACKAGE_NAME
export GO111MODULE=on
go mod vendor
export GO111MODULE=""
cd ../fabric-samples/test-network
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
peer lifecycle chaincode package $PACKAGE_NAME.tar.gz --path ../../chaincode/ --lang golang --label $PACKAGE_NAME

pp_info "Deploy" "Install the chaincode package on Org1"
peer lifecycle chaincode install $PACKAGE_NAME.tar.gz

pp_info "Deploy" "Approve chaincode definition in Org1"
peer lifecycle chaincode queryinstalled
echo
echo "Copy 'Package ID' and past it here:"
read input
export CC_PACKAGE_ID=$input
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID $CHANNEL_NAME --name $PACKAGE_NAME --version 1.0 --signature-policy "OR('Org1MSP.member','Org2MSP.member')" --init-required --package-id $CC_PACKAGE_ID --sequence 1 --tls true --cafile $ORDERER_CA

pp_info "Deploy" "Approve chaincode definition in Org2"
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID $CHANNEL_NAME --name $PACKAGE_NAME --version 1.0 --signature-policy "OR('Org1MSP.member','Org2MSP.member')" --init-required --sequence 1 --tls true --cafile $ORDERER_CA

pp_info "Deploy" "Committing the chaincode definition to the channel"
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORG2_CA=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
sleep 5
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID $CHANNEL_NAME --name $PACKAGE_NAME --version 1.0 --sequence 1 --signature-policy "OR('Org1MSP.member','Org2MSP.member')" --init-required --tls true --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_CA --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_CA

pp_info "Deploy" "Invoking the chaincode (Instantiate)"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID $CHANNEL_NAME --name $PACKAGE_NAME --isInit --tls true --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_CA -c '{"Args":[""]}'

cd ../../bin/deploy-chaincode/
