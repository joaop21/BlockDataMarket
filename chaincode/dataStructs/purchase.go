package dataStructs

import (
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