package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type IdentificationContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *IdentificationContract) Instantiate(_ contractapi.TransactionContextInterface) error {
    return nil
}

// Adds a new Identification to be sell, to the world state with given details
func (_ *IdentificationContract) MakeIdentification(ctx contractapi.TransactionContextInterface, name string, ip string, publicKey string) error {
	
	// create a new Identification
    identification := Identification{
		Type:        "Identification"
		Id: 	     uuid.New().String(),
        Name:        name,
        Ip:          ip,
        PublicKey:   publicKey,
	}
	
	identificationAsBytes, _ := identification.Serialize()
	key, _ := ctx.GetStub().CreateCompositeKey("Identification", []string{
		identification.Id,
	})

	// test if Identification already exists
	obj, _ := ctx.GetStub().GetState(key)
	if obj != nil {
		return fmt.Errorf("identification already exists")
	}

    return ctx.GetStub().PutState(key, identificationAsBytes)
}

// Get all existing Identification on world state 
func (_ *IdentificationContract) GetIdentification(ctx contractapi.TransactionContextInterface, id string) (*Identification, error) {

	key, _ := ctx.GetStub().CreateCompositeKey("Identification", []string{
		id,
	})
	identificationAsBytes, err := ctx.GetStub().GetState(key)
	if identificationAsBytes == nil || err != nil {
		return nil, err
	}

	identification := new (Identification)
	err = identification.Deserialize(identificationAsBytes)
        if err != nil {
            return nil, err
        }

    return identification, nil
}
