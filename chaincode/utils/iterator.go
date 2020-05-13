package utils

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"reflect"
)

// loop an iterator
func GetIteratorValues(resultsIterator shim.StateQueryIteratorInterface, obj interface{}) (res []interface{}, err error)  {
	defer resultsIterator.Close()
	
	for resultsIterator.HasNext(){
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newObj := reflect.New(reflect.TypeOf(obj)).Elem().Interface()
		err = Deserialize(element.Value, newObj)
		if err != nil {
			return nil, err
		}

		res = append(res, newObj)
	}
	return res, err
}
