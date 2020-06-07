package contracts

import (
	"dataMarket/dataStructs"
	"dataMarket/utils"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CategoryContract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (_ *CategoryContract) Instantiate(_ contractapi.TransactionContextInterface) error {
	return nil
}

// Adds a new Category for Announcements to use
func (_ *CategoryContract) MakeCategory(ctx contractapi.TransactionContextInterface, name string, queries []string) (*dataStructs.Category, error) {

	// create a composite key
	key, _ := ctx.GetStub().CreateCompositeKey("Category", []string{name})

	// test if key already exists
	obj, _ := ctx.GetStub().GetState(key)
	if obj != nil {
		return nil, errors.New("category name already exists")
	}

	// remove equal queries
	uniqueQueries := utils.RemoveRepetitions(queries)

	category := dataStructs.NewCategory(name, uniqueQueries)
	if category == nil {
		return nil, errors.New("error creating announcement")
	}

	// serialize category
	categoryAsBytes, _ := utils.Serialize(category)

	// put category in world state
	err := ctx.GetStub().PutState(key, categoryAsBytes)
	if err != nil {
		return nil, errors.New("error putting category in world state")
	}

	return category, nil
}

// Get Category from world state by name
func (_ *CategoryContract) GetCategory(ctx contractapi.TransactionContextInterface, categoryName string) (*dataStructs.Category, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"Category\",\"name\":\"%s\"}}", categoryName)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	results, err := utils.GetIteratorValues(resultsIterator, new(dataStructs.Category))
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("category doesn't exists")
	}

	return results[0].(*dataStructs.Category), nil
}

// Get all Categories from world state
func (_ *CategoryContract) GetCategories(ctx contractapi.TransactionContextInterface) ([]*dataStructs.Category, error) {

	// get all the keys that match with args
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Category", []string{})
	if err != nil {
		return nil, err
	}

	return getCategories(resultsIterator)
}

// Auxiliary function for repeating code
func getCategories(resultsIterator shim.StateQueryIteratorInterface) ([]*dataStructs.Category, error) {

	// Iterate values received
	values, err := utils.GetIteratorValues(resultsIterator, new(dataStructs.Category))
	if err != nil {
		return nil, err
	}

	// Convert to a []Announcement
	categories := convertToCategory(values)
	return categories, err
}

// Converter of an []interface{} to []Announcement
func convertToCategory(values []interface{}) (categories []*dataStructs.Category) {
	categories = make([]*dataStructs.Category, len(values))
	for i := range values {categories[i] = values[i].(*dataStructs.Category)}
	return categories
}




