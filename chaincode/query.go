package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// Query object that represents a query in the World State
type Query struct {
	Type           string    `json:"type"`
	QueryId        string    `json:"queryId"`
	AnnouncementId string    `json:"announcementId"`
	IssuerId       string    `json:"issuerId"`
	Query          string    `json:"query"`
	Price          float32   `json:"price"`
	Response       string    `json:"Response"`
	InsertedAt     time.Time `json:"insertedAt"`
}

// Serialize formats the Query as JSON bytes
func (q *Query) Serialize() ([]byte, error) {
	return json.Marshal(q)
}

// Deserialize formats the Query from JSON bytes
func (q *Query) Deserialize(bytes []byte) error {
	err := json.Unmarshal(bytes, q)

	if err != nil {
		return fmt.Errorf("error deserializing Query. %s", err.Error())
	}

	return nil
}

func GetIteratorValues(resultsIterator shim.StateQueryIteratorInterface) ([]Query, error) {
	defer resultsIterator.Close()

	var res []Query
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newAnn := new(Query)
		err = newAnn.Deserialize(element.Value)
		if err != nil {
			return nil, err
		}

		res = append(res, *newAnn)
	}
	return res, nil
}
