package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Announcement object that represents an announcement in the World State
type Announcement struct {
    AnnouncementId string    `json:"announcementId"`
    DataId         string    `json:"dataId"`
    OwnerId        string    `json:"ownerId"`
    Value          float32   `json:"value"`
    DataCategory   string    `json:"dataCategory"`
    InsertedAt     time.Time `json:"insertedAt"`
}

// Serialize formats the Announcement as JSON bytes
func (ann *Announcement) Serialize() ([]byte, error) {
	return json.Marshal(ann)
}

// Deserialize formats the Announcement from JSON bytes
func Deserialize(bytes []byte, ann *Announcement) error {
	err := json.Unmarshal(bytes, ann)

	if err != nil {
		return fmt.Errorf("error deserializing Announcement. %s", err.Error())
	}

	return nil
}
