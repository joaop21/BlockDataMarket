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
func (ann *Identification) Serialize() ([]byte, error) {
	return json.Marshal(ann)
}

// Deserialize formats the Identification from JSON bytes
func (ann *Identification) Deserialize(bytes []byte, iden *Identification) error {
	err := json.Unmarshal(bytes, ann)

	if err != nil {
		return fmt.Errorf("error deserializing Identification. %s", err.Error())
	}

	return nil
}
