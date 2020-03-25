package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type PurchaseContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *PurchaseContract) Instantiate(_ contractapi.TransactionContextInterface) error {
	return nil
}

// Adds a new Announcement to be sell, to the world state with given details
func (_ *PurchaseContract) MakeAnnouncement(ctx contractapi.TransactionContextInterface,
	dataId string, ownerId string, value float32, cat string) error {

	// check if category is available
	category, err := checkExistence(cat)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// ##### ATTENTION #####
	// check if ownerID exists
	// check if ownerID and the invoking entity are the same
	// Done by the API

	// create a new Announcement
	purchase := PurchaseContract{
		AnnouncementId: uuid.New().String(),
		BuyerId:        buyerId,
		Value:          value,
		InsertedAt:     time.Now(),
	}

	// create a composite key
	purchaseAsBytes, _ := purchase.Serialize()
	key, _ := ctx.GetStub().CreateCompositeKey("Purchase", []string{
		purchase.Announcement,
		purchase.BuyerId,
	})

	// test if key already exists
	obj, _ := ctx.GetStub().GetState(key)
	if obj != nil {
		return fmt.Errorf("key already exists")
	}

	return ctx.GetStub().PutState(key, purchaseAsBytes)
}

// Get all a specific purchase from the world state
func (_ *PurchaseContract) GetPurchase(ctx contractapi.TransactionContextInterface, announcementId string, buyerId string) (Purchase, error) {
	key, _ := ctx.GetStub().CreateCompositeKey("Purchase", []string{
		announcementId,
		buyerId,
	})
	purchaseAsBytes, err := ctx.GetStub().GetState(key)
	if purchaseAsBytes != nil || err != nil {
		return nil, err
	}

	purchase := new Purchase() 
	err = Deserialize(purchaseAsBytes, purchase)
        if err != nil {
            return nil, err
        }

    return purchase, nil
}

// Get all existing Purchases from one Announcement on world state that match with the arguments
func (_ *PurchaseContract) GetAnnouncementPurchases(ctx contractapi.TransactionContextInterface, announcementId string) ([]Purchase, error) {

	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{
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
		err = Deserialize(element.Value, newPurch)
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
		err = Deserialize(element.Value, newPurch)
		if err != nil {
			return nil, err
		}		results = append(results, *newPurch)

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

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return nil, err
	}
	
	return queryResults, nil
}