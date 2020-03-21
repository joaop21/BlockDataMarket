package main

import (
    "github.com/google/uuid"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "time"
)

type AnnouncementContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *AnnouncementContract) Instantiate(_ contractapi.TransactionContextInterface) error {
    return nil
}

// AnnounceData adds a new data record to be sell to the world state with given details
func (_ *AnnouncementContract) MakeAnnouncement(ctx contractapi.TransactionContextInterface,
    dataId string, ownerId string, value float32, category Category) error {

    announcement := Announcement{
        AnnouncementId: uuid.New().String(),
        DataId:         dataId,
        OwnerId:        ownerId,
        Value:          value,
        DataCategory:   category,
        InsertedAt:     time.Now(),
    }

    announcementAsBytes, _ := announcement.Serialize()
    key, _ := ctx.GetStub().CreateCompositeKey("Announcement", []string{
        announcement.DataCategory.String(),
        announcement.OwnerId,
        announcement.AnnouncementId,
    })

    return ctx.GetStub().PutState(key, announcementAsBytes)
}
