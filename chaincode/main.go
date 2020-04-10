package main

import (
	"dataMarket/contracts"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	announcementContract := new(contracts.AnnouncementContract)
	announcementContract.Name = "AnnouncementContract"

	identificationContract := new(contracts.IdentificationContract)
	identificationContract.Name = "IdentificationContract"

	purchaseContract := new(contracts.PurchaseContract)
	purchaseContract.Name = "PurchaseContract"

	queryContract := new(contracts.QueryContract)
	queryContract.Name = "QueryContract"

	chaincode, err := contractapi.NewChaincode(announcementContract, identificationContract, purchaseContract, queryContract)

	if err != nil {
		panic(err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic(err.Error())
	}
}
