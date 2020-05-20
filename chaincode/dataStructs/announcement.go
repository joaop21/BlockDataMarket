package dataStructs

import (
	"time"
)

// Announcement object that represents an announcement in the World State
type Announcement struct {
	Type		     string				`json:"type"`
	AnnouncementId   string    			`json:"announcementId"`
	DataId           string    			`json:"dataId"`
	OwnerId          string    			`json:"ownerId"`
	PossibleQueries  []string      		`json:"possibleQueries"`
	QueryPrices      []float32   		`json:"prices"`
	DataCategory     string    			`json:"dataCategory"`
	InsertedAt       time.Time 			`json:"insertedAt"`
}

// Constructor for Announcement
func NewAnnouncement(announcementId string, dataId string, ownerId string, queries []string, queryPrices []float32, categoryName string, insertionDate time.Time) *Announcement {

	return &Announcement{
		Type: "Announcement",
		AnnouncementId:  announcementId,
		DataId:          dataId,
		OwnerId:         ownerId,
		PossibleQueries: queries,
		QueryPrices:     queryPrices,
		DataCategory:    categoryName,
		InsertedAt:      insertionDate,
	}

}