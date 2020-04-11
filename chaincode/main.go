package main

import (
	"dataMarket/context"
	"dataMarket/contracts"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	announcementContract := new(contracts.AnnouncementContract)
	announcementContract.Name = "AnnouncementContract"
	announcementContract.TransactionContextHandler = new(context.TransactionContext)
	announcementContract.BeforeTransaction = context.SearchIdentitiesHandler

	identificationContract := new(contracts.IdentificationContract)
	identificationContract.Name = "IdentificationContract"
	identificationContract.TransactionContextHandler = new(context.TransactionContext)
	identificationContract.BeforeTransaction = context.SearchIdentitiesHandler

	//purchaseContract := new(contracts.PurchaseContract)
	//purchaseContract.Name = "PurchaseContract"

	queryContract := new(contracts.QueryContract)
	queryContract.Name = "QueryContract"
	queryContract.TransactionContextHandler = new(context.TransactionContext)
	queryContract.BeforeTransaction = context.SearchIdentitiesHandler

	chaincode, err := contractapi.NewChaincode(announcementContract, identificationContract, queryContract)

	if err != nil {
		panic(err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic(err.Error())
	}
}
