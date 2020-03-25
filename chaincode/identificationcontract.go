package main

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "time"
)

type IdentificationContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *IdentificationContract) Instantiate(_ contractapi.TransactionContextInterface) error {
    return nil
}

// Adds a new Identification to be sell, to the world state with given details
func (_ *IdentificationContract) MakeIdentification(ctx contractapi.TransactionContextInterface,
    id string, name string, ip string, publicKey string) error {

	// test if Identification already exists
    obj, _ := ctx.GetStub().GetState(key)
    if obj != nil {
        return fmt.Errorf("Identification already exists")
    }
	
	// create a new Identification
    identification := Identification{
        AnnouncementId: uuid.New().String(),
        Id:          dataId,
        Name:        ownerId,
        Ip:          value,
        PublicKey:   category,
	}
	
	identificationAsBytes, _ := identification.Serialize()
	key, _ := ctx.GetStub().CreateCompositeKey("Identification", []string{
		identification.Id,
	})


    return ctx.GetStub().PutState(key, identificationAsBytes)
}

// Get all existing Identification on world state 
func (_ *IdentificationContract) GetIdentification(ctx contractapi.TransactionContextInterface, string id) (Identification, error) {

	key, _ := ctx.GetStub().CreateCompositeKey("Identification", []string{
		id,
	})
	identificationAsBytes, err := ctx.GetStub().GetState(key)
	if identificationAsBytes != nil || err != nil {
		return nil, err
	}

	identification := new Identification() 
	err = Deserialize(identificationAsBytes, identification)
        if err != nil {
            return nil, err
        }

    return identification, nil
}
