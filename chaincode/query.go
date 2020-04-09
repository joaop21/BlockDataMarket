package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"time"
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

func (q *Query) GetIteratorValues(resultsIterator shim.StateQueryIteratorInterface) ([]Query, error) {
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
