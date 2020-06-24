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
	"strconv"
	"time"
)

type AnnouncementContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *AnnouncementContract) Instantiate(_ context.TransactionContextInterface) error {
	return nil
}

// Adds a new Announcement to be sell, to the world state with given details
func (_ *AnnouncementContract) MakeAnnouncement(ctx context.TransactionContextInterface, dataId string, queries []string, prices []float32, categoryName string) (*dataStructs.Announcement, error) {

	identification := ctx.GetIdentification()
	if identification == nil {
		return nil, errors.New("the submitter has no identification")
	}

	// check if queries and prices have the same length
	if len(queries) != len(prices) {
		return nil, errors.New("queries and correspondent prices have not the same length")
	}

	// check if category exist
	category, err := new(CategoryContract).GetCategory(ctx, categoryName)
	if err != nil {
		return nil, err
	}

	// check if queries are possible
	for _, query := range queries {
		if !utils.Contains(category.PossibleQueries, query) {
			return nil, errors.New("query does not exist in category: " + query)
		}
	}

	// create a new Announcement
	announcement := dataStructs.NewAnnouncement(uuid.New().String(), dataId, ctx.GetIdentification().Id, queries, prices, categoryName, time.Now())

	if announcement == nil {
		return nil, errors.New("error creating announcement")
	}

	// create a composite key
	announcementAsBytes, _ := utils.Serialize(announcement)
	key, _ := ctx.GetStub().CreateCompositeKey("Announcement", []string{
		announcement.DataCategory,
		announcement.OwnerId,
		announcement.AnnouncementId,
	})

	// test if key already exists
	obj, _ := ctx.GetStub().GetState(key)
	if obj != nil {
		return nil, errors.New("key already exists")
	}

	err = ctx.GetStub().PutState(key, announcementAsBytes)
	if err != nil {
		return nil, errors.New("error putting announcement in world state")
	}

	return announcement, nil
}

// Get Announcement on world state by id
func (_ *AnnouncementContract) GetAnnouncement(ctx context.TransactionContextInterface, announcementId string) (*dataStructs.Announcement, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Announcement\",\"announcementId\":\"%s\"}}", announcementId)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	results, err := utils.GetIteratorValues(resultsIterator, new(dataStructs.Announcement))
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("announcement doesn't exists")
	}

	return results[0].(*dataStructs.Announcement), nil
}

// Get all existing Announcements on world state
func (_ *AnnouncementContract) GetAnnouncements(ctx context.TransactionContextInterface) ([]*dataStructs.Announcement, error) {
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{})
	if err != nil {
		return nil, err
	}
	return getAnnouncements(resultsIterator)
}

// Get all Announcements for a category
func (_ *AnnouncementContract) GetAnnouncementsByCategory(ctx context.TransactionContextInterface, categoryName string) ([]*dataStructs.Announcement, error) {
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{categoryName})
	if err != nil {
		return nil, err
	}
	return getAnnouncements(resultsIterator)
}

// Get Announcements for a category lower than a value
func (_ *AnnouncementContract) GetAnnouncementsByCategoryLowerThan(ctx context.TransactionContextInterface, categoryName string, value float32) ([]*dataStructs.Announcement, error) {
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{categoryName})
	if err != nil {
		return nil, err
	}
	announcements, err2 := getAnnouncements(resultsIterator)

	result := []*dataStructs.Announcement{}
	for i := range announcements {
		if hasValidValues(announcements[i], value){
			result = append(result, announcements[i])
		}
	}
	return result, err2
}

// Get all Announcements for an owner
func (_ *AnnouncementContract) GetAnnouncementsByOwner(ctx context.TransactionContextInterface, ownerId string) ([]*dataStructs.Announcement, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"ownerId\":\"%s\"}}", ownerId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return getAnnouncements(resultsIterator)
}

// Get all Announcements lower than a value
func (_ *AnnouncementContract) GetAnnouncementsLowerThan(ctx context.TransactionContextInterface, value float32) ([]*dataStructs.Announcement, error) {
	queryString := fmt.Sprintf("{\"selector\": {\"prices\": {\"$elemMatch\": {\"$lte\": %f}}}}", value)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return getAnnouncements(resultsIterator)
}

// Update an Announcement with new prices
func (_ *AnnouncementContract) UpdateQueryPrices(ctx context.TransactionContextInterface, announcementId string, updates [][]string) (map[string]float32, error) {

	// to change something has to have an identification
	identification := ctx.GetIdentification()
	if identification == nil {
		return nil, errors.New("the submitter has no identification")
	}

	// get the announcement
	announcement, err := new(AnnouncementContract).GetAnnouncement(ctx, announcementId)
	if err != nil {
		return nil, err
	}

	// check the id of the changer
	if announcement.OwnerId != ctx.GetIdentification().Id {
		return nil, errors.New("the changer is not the owner of the announcement")
	}

	// update prices
	lengthPQ := len(announcement.PossibleQueries)
	for _, pair := range updates {

		if len(pair) != 2 {
			return nil, errors.New("bad input")
		}

		price, err := strconv.ParseFloat(pair[1], 32)
		if err != nil {
			return nil, errors.New("bad input, price isn't float:  " + err.Error())
		}

		has := false
		for i := 0 ; i < lengthPQ; i++ {
			if pair[0] == announcement.PossibleQueries[i] {
				announcement.QueryPrices[i] = float32(price)
				has = true
				break
			}
		}

		if !has {
			return nil, errors.New("query doesn't exist in the list of possible queries")
		}

	}

	// create a composite key
	announcementAsBytes, _ := utils.Serialize(announcement)
	key, _ := ctx.GetStub().CreateCompositeKey("Announcement", []string{
		announcement.DataCategory,
		announcement.OwnerId,
		announcement.AnnouncementId,
	})

	// put in the world state
	err = ctx.GetStub().PutState(key, announcementAsBytes)
	if err != nil {
		return nil, errors.New("error putting announcement in world state")
	}

	//build result
	var result map[string]float32
	result = make(map[string]float32)
	for i := 0 ; i < lengthPQ ; i++ {
		result[announcement.PossibleQueries[i]] = announcement.QueryPrices[i]
	}

	return result, nil
}

// Auxiliary function for repeating code
func getAnnouncements(resultsIterator shim.StateQueryIteratorInterface) ([]*dataStructs.Announcement, error) {
	// Iterate values received
	values, err := utils.GetIteratorValues(resultsIterator, new(dataStructs.Announcement))
	if err != nil {
		return nil, err
	}
	// Convert to a []Announcement
	announcements := convertToAnnouncement(values)
	return announcements, err
}

// Converter of an []interface{} to []Announcement
func convertToAnnouncement(values []interface{}) (announcements []*dataStructs.Announcement) {
	announcements = make([]*dataStructs.Announcement, len(values))
	for i := range values {announcements[i] = values[i].(*dataStructs.Announcement)}
	return announcements
}

func hasValidValues(announcement *dataStructs.Announcement, value float32) bool {
	prices := announcement.QueryPrices
	for i := range prices {
		if prices[i] <= value {
			return true
		}
	}

	return false
}

