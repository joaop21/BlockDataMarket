package context

import (
	"dataMarket/dataStructs"
	"dataMarket/utils"
	"errors"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// Function for searching identities before the transaction happens
func SearchIdentitiesHandler(ctx TransactionContextInterface) error {
	// New returns an instance of ClientID
	cliId, err := cid.New(ctx.GetStub())
	if err != nil {
		return errors.New("error initializing client identity")
	}
	// GetID returns a unique ID associated with the invoking identity.
	id, err := cliId.GetID()
	if err != nil {
		return errors.New("error getting ID from client identity")
	}
	// set unique identity in transaction context
	ctx.SetUniqueIdentity(id)
	// create key for search identification
	key, err := ctx.GetStub().CreateCompositeKey("Identification", []string{id})
	if err != nil {
		return errors.New("error creating composite key")
	}
	// get state from the previously generated key
	identificationAsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return errors.New("error getting state")
	}
	// set identification in transaction context
	if identificationAsBytes != nil {
		identification := new(dataStructs.Identification)
		err = utils.Deserialize(identificationAsBytes, identification)
		if err != nil {
			return errors.New("error deserializing identification")
		}
		ctx.SetIdentification(identification)
	} else {
		ctx.SetIdentification(nil)
	}
	return nil
}
