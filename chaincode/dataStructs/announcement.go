package dataStructs

import (
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
	category, err := CheckExistence(categoryName)

	if err != nil || len(queryPrices) != len(category.Actions) {
		return nil
	}

	return &Announcement{
		Type: "Announcement",
		AnnouncementId: announcementId,
		DataId:         dataId,
		OwnerId:        ownerId,
		QueryPrices:    queryPrices,
		DataCategory:   category.Name,
		InsertedAt:     insertionDate,
	}
}