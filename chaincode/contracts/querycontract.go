package contracts

import (
	"dataMarket/context"
	"dataMarket/dataStructs"
	"dataMarket/utils"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type QueryContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *QueryContract) Instantiate(_ context.TransactionContextInterface) error {
	return nil
}

// Adds a new Query to world state
func (_ *QueryContract) MakeQuery(ctx context.TransactionContextInterface, announcementId string, queryArg string, price float32) (*dataStructs.Query, error) {

	identification := ctx.GetIdentification()
	if identification == nil {
		return nil, errors.New("the submitter has no identification")
	}

	// check if announcement exist
	announcement, err := new(AnnouncementContract).GetAnnouncement(ctx, announcementId)
	if err != nil || announcement == nil {
		return nil, errors.New("announcement ID does not exist")
	}

	if !utils.Contains(announcement.PossibleQueries, queryArg) {
		return nil, errors.New("announcement does not support that kind of query")
	}

	// create a new Announcement
	query := dataStructs.NewQuery(announcementId, ctx.GetIdentification().Id, price, queryArg)

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
		return nil, fmt.Errorf("key already exists")
	}

	// send event with query information in payload
	eventName := utils.Concat("Query:", announcementId)
	err = ctx.GetStub().SetEvent(eventName, queryAsBytes)
	if err != nil {
		return nil, errors.New("event can't be emitted")
	}

	err = ctx.GetStub().PutState(key, queryAsBytes)
	if err != nil {
		return nil, errors.New("error putting query in world state")
	}

	return query, nil
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

	eventName := utils.Concat("Response:", queryId)
	err = ctx.GetStub().SetEvent(eventName, queryAsBytes)
	if err != nil {
		return errors.New("event can't be emitted")
	}

	return ctx.GetStub().PutState(key, queryAsBytes)
}

// Get queries made to an announcement
func (_ *QueryContract) GetQueriesByAnnouncement(ctx context.TransactionContextInterface, announcementId string) ([]*dataStructs.Query, error) {
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Query", []string{announcementId})
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

//Gets reponse from specific query
func (_ *QueryContract) GetResponse(ctx context.TransactionContextInterface, queryId string) (string, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Query\",\"queryId\":\"%s\"}}", queryId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "", err
	}
	results, err := utils.GetIteratorValues(resultsIterator, new(dataStructs.Query))
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("Query doesn't exists")
	}

	return results[0].(*dataStructs.Query).Response, nil
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
