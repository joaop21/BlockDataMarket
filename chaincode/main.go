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

	chaincode_announc, err_announc := contractapi.NewChaincode(announcementContract)
	chaincode_ident, err_ident := contractapi.NewChaincode(identificationContract)
	chaincode_purch, err_purch := contractapi.NewChaincode(purchaseContract)

	if err_announc != nil {
		panic(err.Error())
	}

	if err := chaincode_announc.Start(); err != nil {
		panic(err.Error())
	}

	if err_ident != nil {
		panic(err.Error())
	}

	if err := chaincode_ident.Start(); err != nil {
		panic(err.Error())
	}

	if err_purch != nil {
		panic(err.Error())
	}

	if err := chaincode_purch.Start(); err != nil {
		panic(err.Error())
	}
}
