package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Purchase object that represents an purchase in the World State
type Purchase struct {
	Type           string    `json:"type"`
	AnnouncementId string    `json:"announcementId"`
	BuyerId        string    `json:"buyerId"`
	Value          float32   `json:"value"`
	InsertedAt     time.Time `json:"insertedAt"`
}

// Serialize formats the Purchase as JSON bytes
func (purch *Purchase) Serialize() ([]byte, error) {
	return json.Marshal(purch)
}

// Deserialize formats the Purchase from JSON bytes
func (purch *Purchase) Deserialize(bytes []byte) error {
	err := json.Unmarshal(bytes, purch)

	if err != nil {
		return fmt.Errorf("error deserializing Purchase. %s", err.Error())
	}

	return nil
}
