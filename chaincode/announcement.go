package main

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
	}
}

// Converter of an []interface{} to []Announcement
func ConvertToAnnouncement(values []interface{}) (announcements []Announcement) {
	announcements = make([]Announcement, len(values))
	for i := range values {
		announcements[i] = values[i].(Announcement)
	}
	return announcements
}
