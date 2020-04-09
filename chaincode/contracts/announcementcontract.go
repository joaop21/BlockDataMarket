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

type AnnouncementContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *AnnouncementContract) Instantiate(_ contractapi.TransactionContextInterface) error {
	return nil
}

// Adds a new Announcement to be sell, to the world state with given details
func (_ *AnnouncementContract) MakeAnnouncement(ctx contractapi.TransactionContextInterface, dataId string, ownerId string, prices []float32, categoryName string) error {

	// create a new Announcement
	announcement := dataStructs.NewAnnouncement(uuid.New().String(), dataId, ownerId, prices, categoryName, time.Now())

	if announcement == nil {
		return fmt.Errorf("Error creating announcement")
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
		return fmt.Errorf("key already exists")
	}

	return ctx.GetStub().PutState(key, announcementAsBytes)
}

// Get Announcement on world state by id
func (_ *AnnouncementContract) GetAnnouncement(ctx contractapi.TransactionContextInterface, announcementId string) (*dataStructs.Announcement, error) {

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
		return nil, fmt.Errorf("Announcement doesn't exists")
	}

	return results[0].(*dataStructs.Announcement), nil
}

// Get all existing Announcements on world state
func (_ *AnnouncementContract) GetAnnouncements(ctx contractapi.TransactionContextInterface) ([]*dataStructs.Announcement, error) {
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{})
	if err != nil {
		return nil, err
	}
	return getAnnouncements(resultsIterator)
}

// Get all Announcements for a category
func (_ *AnnouncementContract) GetAnnouncementsByCategory(ctx contractapi.TransactionContextInterface, categoryName string) ([]*dataStructs.Announcement, error) {
	// check if category is available
	category, err := dataStructs.CheckExistence(categoryName)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{category.Name})
	if err != nil {
		return nil, err
	}
	return getAnnouncements(resultsIterator)
}

// Get all Announcements for an owner
func (_ *AnnouncementContract) GetAnnouncementsByOwner(ctx contractapi.TransactionContextInterface, ownerId string) ([]*dataStructs.Announcement, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"ownerId\":\"%s\"}}", ownerId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return getAnnouncements(resultsIterator)
}

// Get all Announcements lower than a value
func (_ *AnnouncementContract) GetAnnouncementsLowerThan(ctx contractapi.TransactionContextInterface, value float32) ([]*dataStructs.Announcement, error) {
	queryString := fmt.Sprintf("{\"selector\": {\"prices\": {\"$elemMatch\": {\"$lte\": %f}}}}", value)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return getAnnouncements(resultsIterator)
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

