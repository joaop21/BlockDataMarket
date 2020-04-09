package main

import (
	"dataMarket/utils"
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
func (_ *QueryContract) MakeQuery(ctx contractapi.TransactionContextInterface, announcementId string, issuerId string, queryArg string, price float32) error {

	// create a new Announcement
	query := Query{
		Type:			"Query",
		QueryId:        uuid.New().String(),
		AnnouncementId: announcementId,
		IssuerId:       issuerId,
		Price:			price,
		Query:          queryArg,
		Response:   	"",
		InsertedAt:     time.Now(),
	}

	// create a composite key
	queryAsBytes, _ := utils.Serialize(query)
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
func (_ *QueryContract) PutResponse(ctx contractapi.TransactionContextInterface, queryid string, response string) error {

	var results []Query
	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Query\",\"queryId\":\"%s\"}}", queryid)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return err
	}
	query := new(Query)
	values, err := utils.GetIteratorValues(resultsIterator)
	results = ConvertToQuery(values)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return fmt.Errorf("Query doesn't exists")
	}

	query = &(results[0])
	query.Response = response
	var queryAsBytes []byte
	queryAsBytes, _ = utils.Serialize(query)

	key, _ := ctx.GetStub().CreateCompositeKey("Query", []string{
                query.AnnouncementId,
                query.IssuerId,
                query.QueryId,
        })


	return ctx.GetStub().PutState(key, queryAsBytes)
}

// Get queries made to an announcement
func (_ *QueryContract) GetQueriesByAnnouncement(ctx contractapi.TransactionContextInterface, announcementId string) ([]Query, error) {

	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Query", []string{announcementId,})
	if err != nil {
		return nil, err
	}
	values, err := utils.GetIteratorValues(resultsIterator)
	res := ConvertToQuery(values)
	return res, nil
}

// Get query by its id
func (_ *QueryContract) GetQuery(ctx contractapi.TransactionContextInterface,
	queryId string) (*Query, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Query\",\"queryId\":\"%s\"}}", queryId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	results, err := utils.GetIteratorValues(resultsIterator)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("Query doesn't exists")
	}

	return results[0].(*Query), nil
}


// Get queries made to an announcement by an issuer
func (_ *QueryContract) GetQueriesByIssuer(ctx contractapi.TransactionContextInterface, issuerId string) ([]Query, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"issuerId\":\"%s\"}}", issuerId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	values, err := utils.GetIteratorValues(resultsIterator)
	res := ConvertToQuery(values)
	return res, nil
}
