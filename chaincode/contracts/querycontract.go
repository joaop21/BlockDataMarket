package contracts

import (
	"dataMarket/dataStructs"
	"dataMarket/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
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
	query := dataStructs.Query{
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

	var results []*dataStructs.Query
	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Query\",\"queryId\":\"%s\"}}", queryid)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return err
	}
	query := new(dataStructs.Query)
	results, err = getQueries(resultsIterator)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return fmt.Errorf("Query doesn't exists")
	}

	query = results[0]
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
func (_ *QueryContract) GetQueriesByAnnouncement(ctx contractapi.TransactionContextInterface, announcementId string) ([]*dataStructs.Query, error) {
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Query", []string{announcementId,})
	if err != nil {
		return nil, err
	}
	return getQueries(resultsIterator)
}

// Get query by its id
func (_ *QueryContract) GetQuery(ctx contractapi.TransactionContextInterface,
	queryId string) (*dataStructs.Query, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Query\",\"queryId\":\"%s\"}}", queryId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	results, err := utils.GetIteratorValues(resultsIterator, new(dataStructs.Query))
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("Query doesn't exists")
	}

	return results[0].(*dataStructs.Query), nil
}


// Get queries made to an announcement by an issuer
func (_ *QueryContract) GetQueriesByIssuer(ctx contractapi.TransactionContextInterface, issuerId string) ([]*dataStructs.Query, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"issuerId\":\"%s\"}}", issuerId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return getQueries(resultsIterator)
}

// Auxiliary function for repeating code
func getQueries(resultsIterator shim.StateQueryIteratorInterface) ([]*dataStructs.Query, error) {
	// Iterate values received
	values, err := utils.GetIteratorValues(resultsIterator, new(dataStructs.Query))
	if err != nil {
		return nil, err
	}
	// Convert to a []Announcement
	queries := convertToQuery(values)
	return queries, err
}

// Converter of an []interface{} to []Query
func convertToQuery(values []interface{}) (queries []*dataStructs.Query) {
	queries = make([]*dataStructs.Query, len(values))
	for i := range values {
		queries[i] = values[i].(*dataStructs.Query)
	}
	return queries
}
