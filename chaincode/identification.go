package main

import (
	"encoding/json"
	"fmt"
)

// Identification object that represents an identification in the World State
type Identification struct {
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
