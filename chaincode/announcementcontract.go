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

    // check if category is available
    category, err := checkExistence(cat)
    if err != nil {
        return fmt.Errorf(err.Error())
    }

    // ##### ATTENTION #####
    // check if ownerID exists
    // check if ownerID and the invoking entity are the same

    // create a new Announcement
    announcement := Announcement{
        AnnouncementId: uuid.New().String(),
        DataId:         dataId,
        OwnerId:        ownerId,
        Value:          value,
        DataCategory:   category,
        InsertedAt:     time.Now(),
    }

    // create a composite key
    announcementAsBytes, _ := announcement.Serialize()
    key, _ := ctx.GetStub().CreateCompositeKey("Announcement", []string{
        announcement.DataCategory,
        announcement.OwnerId,
        announcement.AnnouncementId,
    })

    // test if key already exists
    obj, _ := ctx.GetStub().GetState(key)
    if obj != nil {
        return fmt.Errorf("key already exists")
    }

    return ctx.GetStub().PutState(key, announcementAsBytes)
}

// Get all existing Announcements on world state that match with the arguments
func (_ *AnnouncementContract) GetAnnouncements(ctx contractapi.TransactionContextInterface,
    args ...string) ([]Announcement, error) {

    // get all the keys that match with args
    resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Announcement", args)
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
        announcementAsBytes, err := ctx.GetStub().GetState(element.Key)
        if err != nil {
            return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
        }
        err = Deserialize(announcementAsBytes, new(Announcement))
        if err != nil {
            return nil, err
        }

        res = append(res, *newAnn)
    }

    return res, nil
}
