package utils

import (
	"encoding/json"
	"fmt"
)

// Serialize function for marshal a struct into bytes
func Serialize(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

// Deserialize function for unmarshal bytes into struct
func Deserialize(bytes []byte, i interface{}) error {
	err := json.Unmarshal(bytes, i)
	if err != nil {
		return fmt.Errorf("error deserializing struct. %s", err.Error())
	}
	return nil
}
