package dataStructs

import (
	"github.com/google/uuid"
	"time"
)

// Query object that represents a query in the World State
type Query struct {
	Type           string    `json:"type"`
	QueryId        string    `json:"queryId"`
	AnnouncementId string    `json:"announcementId"`
	IssuerId       string    `json:"issuerId"`
	Query          string    `json:"query"`
	Price          float32   `json:"price"`
	Response       string    `json:"response"`
	InsertedAt     time.Time `json:"insertedAt"`
}

// Constructor for Query
func NewQuery(announcementId string, issuerId string, price float32, queryArg string) *Query {

	return &Query{
		Type:           "Query",
		QueryId:        uuid.New().String(),
		AnnouncementId: announcementId,
		IssuerId:       issuerId,
		Price:          price,
		Query:          queryArg,
		Response:       "",
		InsertedAt:     time.Now(),
	}

}