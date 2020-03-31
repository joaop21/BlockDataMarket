package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Query object that represents a query in the World State
type Query struct {
	QueryId        string     `json:"queryId"`
	AnnouncementId string     `json:"announcementId"`
	IssuerId       string     `json:"issuerId"`
	Query          string     `json:"query"`
	Price		   float32	  `json:"price"`
	Response       string     `json:"Response"`
	InsertedAt     time.Time  `json:"insertedAt"`
}

// Serialize formats the Query as JSON bytes
func (q *Query) Serialize() ([]byte, error) {
	return json.Marshal(q)
}

// Deserialize formats the Query from JSON bytes
func (q *Query) Deserialize(bytes []byte) error {
	err := json.Unmarshal(bytes, q)

	if err != nil {
		return fmt.Errorf("error deserializing Query. %s", err.Error())
	}

	return nil
}