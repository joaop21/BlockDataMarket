package context

import (
	"dataMarket/dataStructs"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Create a new TransactionContextInterface that has contractapi.TransactionContextInterface for composition
// it's like a re-definition
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetIdentification() *dataStructs.Identification
	GetUniqueIdentity() string
	SetIdentification(*dataStructs.Identification)
	SetUniqueIdentity(uniqId string)
}

// Re-defining TransactionContext that has contractapi.TransactionContext for composition
type TransactionContext struct {
	contractapi.TransactionContext
	identification *dataStructs.Identification
	uniqueIdentity string
}

func (ic *TransactionContext) GetIdentification() *dataStructs.Identification {
	return ic.identification
}

func (ic *TransactionContext) SetIdentification(ident *dataStructs.Identification)  {
	ic.identification = ident
}

func (ic *TransactionContext) GetUniqueIdentity() string {
	return ic.uniqueIdentity
}

func (ic *TransactionContext) SetUniqueIdentity(uniqId string)  {
	ic.uniqueIdentity = uniqId
}


