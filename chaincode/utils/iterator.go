package utils

import (
	"errors"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/jinzhu/copier"
)

// loop an iterator
func GetIteratorValues(resultsIterator shim.StateQueryIteratorInterface, obj interface{}) (res []interface{}, err error)  {
	defer resultsIterator.Close()
	
	for resultsIterator.HasNext(){
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newObj := obj
		err = copier.Copy(newObj, obj)
		if err != nil {
			return nil, errors.New("can't deep clone obj")
		}

		err = Deserialize(element.Value, newObj)
		if err != nil {
			return nil, err
		}

		res = append(res, newObj)
	}
	return res, err
}
