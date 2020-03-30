package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

type QueryContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *QueryContract) Instantiate(_ contractapi.TransactionContextInterface) error {
	return nil
}

// Adds a new Query to world state
func (_ *QueryContract) MakeQuery(ctx contractapi.TransactionContextInterface,
	announcementId string, issuerId string, queryArg string, price int) error {

	// create a new Announcement
	query := Query{
		QueryId:        uuid.New().String(),
		AnnouncementId: announcementId,
		IssuerId:       issuerId,
		Price:			price,
		Query:          queryArg,
		Response:   	"",
		InsertedAt:     time.Now(),
	}

	// create a composite key
	queryAsBytes, _ := query.Serialize()
	key, _ := ctx.GetStub().CreateCompositeKey("Query", []string{
		query.AnnouncementId,
		query.IssuerId,
		query.QueryId,
	})

	// test if key already exists
	obj, _ := ctx.GetStub().GetState(key)
	if obj != nil {
		return fmt.Errorf("key already exists")
	}

	return ctx.GetStub().PutState(key, queryAsBytes)
}

// Adds a new Query to world state
func (_ *QueryContract) PutResponse(ctx contractapi.TransactionContextInterface,
	announcementId string, issuerId string, queryid string, response string) error {

	key, _ := ctx.GetStub().CreateCompositeKey("Query", []string{
		announcementId,
		issuerId,
		queryid,
	})

	queryAsBytes, err := ctx.GetStub().GetState(key)
	if queryAsBytes == nil || err != nil {
		return err
	}

	query := new(Query)
	err = query.Deserialize(queryAsBytes)
	if err != nil {
		return err
	}

	query.Response = response
	queryAsBytes, _ = query.Serialize()

	return ctx.GetStub().PutState(key, queryAsBytes)
}

// Get queries made to an announcement
func (_ *QueryContract) GetQueriesByAnnouncement(ctx contractapi.TransactionContextInterface,
	announcementId string) ([]Query, error) {

	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Query", []string{announcementId,})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var res []Query
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newQuery := new(Query)
		err = newQuery.Deserialize(element.Value)
		if err != nil {
			return nil, err
		}

		res = append(res, *newQuery)
	}
	return res, nil
}

// Get queries made to an announcement by an issuer
func (_ *QueryContract) GetQueriesByIssuer(ctx contractapi.TransactionContextInterface,
	issuerId string) ([]Query, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"issuerId\":\"%s\"}}", issuerId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var res []Query
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newQuery := new(Query)
		err = newQuery.Deserialize(element.Value)
		if err != nil {
			return nil, err
		}

		res = append(res, *newQuery)
	}
	return res, nil
}