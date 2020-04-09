package main

import (
	"./utils"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

type PurchaseContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *PurchaseContract) Instantiate(_ contractapi.TransactionContextInterface) error {
	return nil
}

// Adds a new Announcement to be sell, to the world state with given details
func (_ *PurchaseContract) MakePurchase(ctx contractapi.TransactionContextInterface, announcementId string, buyerId string, value float32) error {

	// ##### ATTENTION #####
	// check if ownerID exists
	// check if ownerID and the invoking entity are the same
	// Done by the API

	// create a new Announcement
	purchase := Purchase{
		AnnouncementId: announcementId,
		BuyerId:        buyerId,
		Value:          value,
		InsertedAt:     time.Now(),
	}

	// create a composite key
	purchaseAsBytes, _ := utils.Serialize(purchase)
	key, _ := ctx.GetStub().CreateCompositeKey("Purchase", []string{
		purchase.AnnouncementId,
		purchase.BuyerId,
	})

	// test if key already exists
	obj, _ := ctx.GetStub().GetState(key)
	if obj != nil {
		purch := new (Purchase)
		err := utils.Deserialize(obj, purch)
		if err != nil {
			return err
		}
		// if the new purchase has a lower value than the existent one
		if purch.Value >= purchase.Value {
			return fmt.Errorf("purchase has been made with a higher value")
		}
	}

	return ctx.GetStub().PutState(key, purchaseAsBytes)
}

// Get all a specific purchase from the world state
func (_ *PurchaseContract) GetPurchase(ctx contractapi.TransactionContextInterface, announcementId string, buyerId string) (*Purchase, error) {
	key, _ := ctx.GetStub().CreateCompositeKey("Purchase", []string{
		announcementId,
		buyerId,
	})
	purchaseAsBytes, err := ctx.GetStub().GetState(key)
	if purchaseAsBytes == nil || err != nil {
		return nil, err
	}

	purchase := new (Purchase)
	err = utils.Deserialize(purchaseAsBytes, purchase)
        if err != nil {
            return nil, err
        }

    return purchase, nil
}

// Get all existing Purchases from one Announcement on world state that match with the arguments
func (_ *PurchaseContract) GetAnnouncementPurchases(ctx contractapi.TransactionContextInterface, announcementId string) ([]Purchase, error) {

	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Purchase", []string{
		announcementId,
	})

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var res []Purchase
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newPurch := new(Purchase)
		err = utils.Deserialize(element.Value, newPurch)
		if err != nil {
			return nil, err
		}

		res = append(res, *newPurch)
	}

	return res, nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]Purchase, error) {
	var results []Purchase
	for resultsIterator.HasNext() {
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newPurch := new(Purchase)
		err = utils.Deserialize(element.Value, newPurch)
		if err != nil {
			return nil, err
		}
		results = append(results, *newPurch)
	}

	return results, nil
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]Purchase, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	purchases, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return purchases, nil
}

// Get all existing Purchases from one buyer on world state that match with the arguments
func (_ *PurchaseContract) GetBuyerPurchases(ctx contractapi.TransactionContextInterface, buyerId string) ([]Purchase, error) {
	
	queryString := fmt.Sprintf("{\"selector\":{\"buyerId\":\"%s\"}}", buyerId)

	queryResults, err := getQueryResultForQueryString(ctx.GetStub() , queryString)
	if err != nil {
		return nil, err
	}
	
	return queryResults, nil
}
