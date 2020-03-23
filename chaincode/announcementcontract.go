package main

import (
    "fmt"
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

// Adds a new Announcement to be sell, to the world state with given details
func (_ *AnnouncementContract) MakeAnnouncement(ctx contractapi.TransactionContextInterface,
    dataId string, ownerId string, value float32, cat string) error {

    category, err := checkExistence(cat)

    if err != nil {
        return fmt.Errorf(err.Error())
    }

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
        announcement.DataCategory,
        announcement.OwnerId,
        announcement.AnnouncementId,
    })

    return ctx.GetStub().PutState(key, announcementAsBytes)
}

// Get all existing Announcements on world state
func (_ *AnnouncementContract) GetAnnoucements(ctx contractapi.TransactionContextInterface) ([]Announcement, error) {
    resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", []string{})
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var res []Announcement
    var i int
    for i = 0; resultsIterator.HasNext(); i++ {
        element, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        newAnn := new(Announcement)
        err = Deserialize(element.Value, new(Announcement))
        if err != nil {
            return nil, err
        }

        res = append(res, *newAnn)
    }

    return res, nil
}
