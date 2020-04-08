package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"time"
)

// Announcement object that represents an announcement in the World State
type Announcement struct {
	Type		   string				`json:"type"`
	AnnouncementId string    			`json:"announcementId"`
	DataId         string    			`json:"dataId"`
	OwnerId        string    			`json:"ownerId"`
	QueryPrices    []float32   			`json:"prices"`
	DataCategory   string    			`json:"dataCategory"`
	InsertedAt     time.Time 			`json:"insertedAt"`
}

// Constructor for Announcement
func NewAnnouncement(announcementId string, dataId string, ownerId string, queryPrices []float32, categoryName string, insertionDate time.Time) *Announcement {
	category, err := checkExistence(categoryName)

	if err != nil || len(queryPrices) != len(category.actions) {
		return nil
	}

	return &Announcement{
		Type: "Announcement",
		AnnouncementId: announcementId,
		DataId:         dataId,
		OwnerId:        ownerId,
		QueryPrices:    queryPrices,
		DataCategory:   category.name,
		InsertedAt:     insertionDate,
	};
}

// Serialize formats the Announcement as JSON bytes
func (ann *Announcement) Serialize() ([]byte, error) {
	return json.Marshal(ann)
}

// Deserialize formats the Announcement from JSON bytes
func (ann *Announcement) Deserialize(bytes []byte) error {
	err := json.Unmarshal(bytes, ann)

	if err != nil {
		return fmt.Errorf("error deserializing Announcement. %s", err.Error())
	}

	return nil
}

// loop an iterator
func (ann *Announcement) GetIteratorValues(resultsIterator shim.StateQueryIteratorInterface) ([]Announcement, error)  {
	defer resultsIterator.Close()

	var res []Announcement
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newAnn := new(Announcement)
		err = newAnn.Deserialize(element.Value)
		if err != nil {
			return nil, err
		}

		res = append(res, *newAnn)
	}
	return res, nil
}

