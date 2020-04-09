package main

import (
	"dataMarket/utils"
	"fmt"
	"github.com/google/uuid"
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
	announcement := NewAnnouncement(uuid.New().String(), dataId, ownerId, prices, categoryName, time.Now())

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
func (_ *AnnouncementContract) GetAnnouncement(ctx contractapi.TransactionContextInterface, announcementId string) (*interface{}, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Announcement\",\"announcementId\":\"%s\"}}", announcementId)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	results, err := utils.GetIteratorValues(resultsIterator)
	if err != nil {
		return nil, err
	}	
	if len(results) == 0 {
		return nil, fmt.Errorf("Announcement doesn't exists")
	}

	return &(results[0]), nil
}

// Get all existing Announcements on world state
func (_ *AnnouncementContract) GetAnnouncements(ctx contractapi.TransactionContextInterface) ([]interface{}, error) {
	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{})
	if err != nil {
		return nil, err
	}
	return utils.GetIteratorValues(resultsIterator)
}

// Get all Announcements for a category
func (_ *AnnouncementContract) GetAnnouncementsByCategory(ctx contractapi.TransactionContextInterface, categoryName string) ([]interface{}, error) {

	// check if category is available
	category, err := checkExistence(categoryName)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{category.name})
	if err != nil {
		return nil, err
	}
	return utils.GetIteratorValues(resultsIterator)
}

// Get all Announcements for an owner
func (_ *AnnouncementContract) GetAnnouncementsByOwner(ctx contractapi.TransactionContextInterface, ownerId string) ([]interface{}, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"ownerId\":\"%s\"}}", ownerId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return utils.GetIteratorValues(resultsIterator)
}

// Get all Announcements lower than a value
func (_ *AnnouncementContract) GetAnnouncementsLowerThan(ctx contractapi.TransactionContextInterface, value float32) ([]interface{}, error) {

	queryString := fmt.Sprintf("{\"selector\": {\"prices\": {\"$elemMatch\": {\"$lte\": %f}}}}", value)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return utils.GetIteratorValues(resultsIterator)
}
