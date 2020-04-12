package contracts

import (
	"dataMarket/context"
	"dataMarket/dataStructs"
	"dataMarket/utils"
	"errors"
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
func (_ *QueryContract) Instantiate(_ context.TransactionContextInterface) error {
	return nil
}

// Adds a new Query to world state
func (_ *QueryContract) MakeQuery(ctx context.TransactionContextInterface, announcementId string, queryArg string, price float32) error {

	identification := ctx.GetIdentification()
	if identification == nil {
		return errors.New("the submitter has no identification")
	}

	// create a new Announcement
	query := dataStructs.Query{
		Type:			"Query",
		QueryId:        uuid.New().String(),
		AnnouncementId: announcementId,
		IssuerId:       ctx.GetIdentification().Id,
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
func (_ *QueryContract) PutResponse(ctx context.TransactionContextInterface, queryId string, response string) error {

	identification := ctx.GetIdentification()
	if identification == nil {
		return errors.New("the submitter has no identification")
	}

	var results []*dataStructs.Query
	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Query\",\"queryId\":\"%s\"}}", queryId)
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
	// check the submitter identity
	announcement, err := new(AnnouncementContract).GetAnnouncement(ctx, query.AnnouncementId)
	if err != nil {
		return err
	}
	if announcement.OwnerId != identification.Id {
		return errors.New("the submitter isn't the announcement owner. He can't respond")
	}

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
func (_ *QueryContract) GetQueriesByAnnouncement(ctx context.TransactionContextInterface, announcementId string) ([]*dataStructs.Query, error) {
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Query", []string{announcementId,})
	if err != nil {
		return nil, err
	}
	return getQueries(resultsIterator)
}

// Get query by its id
func (_ *QueryContract) GetQuery(ctx context.TransactionContextInterface, queryId string) (*dataStructs.Query, error) {

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
func (_ *QueryContract) GetQueriesByIssuer(ctx context.TransactionContextInterface, issuerId string) ([]*dataStructs.Query, error) {
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
