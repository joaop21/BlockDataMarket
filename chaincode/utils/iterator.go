package utils

import "github.com/hyperledger/fabric-chaincode-go/shim"

// loop an iterator
func GetIteratorValues(resultsIterator shim.StateQueryIteratorInterface) ([]interface{}, error)  {
	defer resultsIterator.Close()

	var res []interface{}
	for resultsIterator.HasNext(){
		element, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newObj := new(interface{})
		err = Deserialize(element.Value, newObj)
		if err != nil {
			return nil, err
		}

		res = append(res, *newObj)
	}
	return res, nil
}