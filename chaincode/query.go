package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Announcement object that represents an announcement in the World State
type Query struct {
	AnnouncementId string    `json:"announcementId"`
	IssuerId       string    `json:"issuerId"`
	Query          string    `json:"query"`
	Response       string    `json:"Response"`
	InsertedAt     time.Time `json:"insertedAt"`
}

// Serialize formats the Announcement as JSON bytes
func (q *Query) Serialize() ([]byte, error) {
	return json.Marshal(q)
}

// Deserialize formats the Announcement from JSON bytes
func (q *Query) Deserialize(bytes []byte) error {
	err := json.Unmarshal(bytes, q)

	if err != nil {
		return fmt.Errorf("error deserializing Query. %s", err.Error())
	}

	return nil
}