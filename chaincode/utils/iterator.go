package utils

import "github.com/hyperledger/fabric-chaincode-go/shim"

// loop an iterator
func GetIteratorValues(resultsIterator shim.StateQueryIteratorInterface, obj interface{}) (res []interface{}, err error)  {
	defer resultsIterator.Close()
	
	for resultsIterator.HasNext(){
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newObj := obj
		err = Deserialize(element.Value, newObj)
		if err != nil {
			return nil, err
		}

		res = append(res, newObj)
	}
	return res, err
}