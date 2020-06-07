package contracts

import (
	"dataMarket/context"
	"dataMarket/dataStructs"
	"dataMarket/utils"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type IdentificationContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *IdentificationContract) Instantiate(_ context.TransactionContextInterface) error {
    return nil
}

// Adds a new Identification to be sell, to the world state with given details
func (_ *IdentificationContract) MakeIdentification(ctx context.TransactionContextInterface, name string, publicKey string) error {

	if ctx.GetIdentification() != nil {
		return errors.New("submitter already exists")
	}

	// create a new Identification
    identification := dataStructs.Identification{
		Type:        "Identification",
		Id: 	     ctx.GetUniqueIdentity(),
        Name:        name,
        PublicKey:   publicKey,
	}
	
	identificationAsBytes, _ := utils.Serialize(identification)
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
func (_ *IdentificationContract) GetIdentification(ctx context.TransactionContextInterface, id string) (*dataStructs.Identification, error) {

	key, _ := ctx.GetStub().CreateCompositeKey("Identification", []string{
		id,
	})
	identificationAsBytes, err := ctx.GetStub().GetState(key)
	if identificationAsBytes == nil || err != nil {
		return nil, err
	}

	identification := new (dataStructs.Identification)
	err = utils.Deserialize(identificationAsBytes, identification)
	if err != nil {
        return nil, err
    }

    return identification, nil
}
