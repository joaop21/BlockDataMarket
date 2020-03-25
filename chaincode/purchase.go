package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Purchase object that represents an purchase in the World State
type Purchase struct {
	AnnouncementId string    `json:"announcementId"`
	BuyerId        string    `json:"buyerId"`
	Value          float32   `json:"value"`
	InsertedAt     time.Time `json:"insertedAt"`
}

// Serialize formats the Purchase as JSON bytes
func (ann *Purchase) Serialize() ([]byte, error) {
	return json.Marshal(ann)
}

// Deserialize formats the Purchase from JSON bytes
func (ann *Purchase) Deserialize(bytes []byte, ann *Purchase) error {
	err := json.Unmarshal(bytes, ann)

	if err != nil {
		return fmt.Errorf("error deserializing Purchase. %s", err.Error())
	}

	return nil
}
