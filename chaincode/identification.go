package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// Identification object that represents an identification in the World State
type Identification struct {
	Type      string `json:"type"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	PublicKey string `json:"publicKey"`
}

// Serialize formats the Identification as JSON bytes
func (iden *Identification) Serialize() ([]byte, error) {
	return json.Marshal(iden)
}

// Deserialize formats the Identification from JSON bytes
func (iden *Identification) Deserialize(bytes []byte) error {
	err := json.Unmarshal(bytes, iden)

	if err != nil {
		return fmt.Errorf("error deserializing Identification. %s", err.Error())
	}

	return nil
}

func (iden *Identification) GetIteratorValues(resultsIterator shim.StateQueryIteratorInterface) ([]Identification, error) {
	defer resultsIterator.Close()

	var res []Identification
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newAnn := new(Identification)
		err = newAnn.Deserialize(element.Value)
		if err != nil {
			return nil, err
		}

		res = append(res, *newAnn)
	}
	return res, nil
}
