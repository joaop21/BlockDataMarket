package main

import (
    "fmt"
    "encoding/json"
)

type Announcement struct {
    AnnouncementId string   `json:"announcementId"`
    DataId         string   `json:"dataId"`
    OwnerId        string   `json:"ownerId"`
    Value          float32  `json:"value"`
    DataCategory   Category `json:"value"`
}

// Serialize formats the Announcement as JSON bytes
func (ann *Announcement) Serialize() ([]byte, error) {
	return json.Marshal(ann)
}

// Deserialize formats the Announcement from JSON bytes
func Deserialize(bytes []byte, ann *Announcement) error {
	err := json.Unmarshal(bytes, ann)

	if err != nil {
		return fmt.Errorf("Error deserializing commercial paper. %s", err.Error())
	}

	return nil
}