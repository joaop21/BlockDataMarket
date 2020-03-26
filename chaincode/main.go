package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	announcementContract := new(AnnouncementContract)
	announcementContract.Name = "AnnouncementContract"

	identificationContract := new(IdentificationContract)
	identificationContract.Name = "IdentificationContract"

	purchaseContract := new(PurchaseContract)
	purchaseContract.Name = "PurchaseContract"

	queryContract := new(QueryContract)
	queryContract.Name = "QueryContract"

	chaincode, err := contractapi.NewChaincode(announcementContract, identificationContract, purchaseContract, queryContract)

	if err != nil {
		panic(err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic(err.Error())
	}
}
