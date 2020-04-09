package main

import (
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
	Response       string    `json:"Response"`
	InsertedAt     time.Time `json:"insertedAt"`
}

// Converter of an []interface{} to []Query
func ConvertToQuery(values []interface{}) (queries []Query) {
	queries = make([]Query, len(values))
	for i := range values {
		queries[i] = values[i].(Query)
	}
	return queries
}
