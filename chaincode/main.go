package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	announcementContract := new(AnnouncementContract)
	announcementContract.Name = "AnnouncementContract"

	chaincode, err := contractapi.NewChaincode(announcementContract)

	if err != nil {
		panic(err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic(err.Error())
	}
}
