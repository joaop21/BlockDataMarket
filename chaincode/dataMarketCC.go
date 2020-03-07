package main

import (
        "encoding/json"
        "fmt"
	"time"
	"github.com/hyperledger/fabric-chaincode-go/shim"
        "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a market
type SmartContract struct {
        contractapi.Contract
}

// Data describes basic details of what makes up a Data selling record
type Data struct {
    Data_id   string `json:"data_id"`
    Creator_id  string `json:"creator_id"`
    Owner_id  string `json:"owner_id"`
    Value  float32 `json:"value"`
    Inserted_At  time.Time `json:"inserted_at"`
    Category string `json:category`
    
}

// QueryResult structure used for handling result of query
type QueryResult struct {
        Key    string `json:"Key"`
        Record *Data
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
  return nil
}

// AnnounceData adds a new data record to be sell to the world state with given details
func (s *SmartContract) AnnounceData(ctx contractapi.TransactionContextInterface, 
  record_id string, data_id string, creator_id string, value float32) error {
  data := Data{
      Data_id: data_id,
      Creator_id: creator_id,
      Owner_id: "",
      Value: value,
      Inserted_At: time.Now(),
  }

  dataAsBytes, _ := json.Marshal(data)

  return ctx.GetStub().PutState(record_id, dataAsBytes)
}

// BuyData adds a new record selled the world state with given details
func (s *SmartContract) BuyData(ctx contractapi.TransactionContextInterface, 
	record_id string, data Data, owner_id string) error {
  new_data := Data{
      Data_id: data.Data_id,
      Creator_id: data.Creator_id,
      Owner_id: owner_id,
      Value: data.Value,
      Inserted_At: time.Now(),
  }

  dataAsBytes, _ := json.Marshal(new_data)

  return ctx.GetStub().PutState(record_id, dataAsBytes)
}

// GetDataId returns de data id of a data record with given details
func (s *SmartContract) GetDataId(ctx contractapi.TransactionContextInterface, 
	record_id string) (string, error) {
	
  dataAsBytes, err := ctx.GetStub().GetState(record_id)

  if err != nil {
	  return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
  }

  if dataAsBytes == nil {
      return "", fmt.Errorf("%s does not exist", record_id)
  }

  data := new(Data)
  _ = json.Unmarshal(dataAsBytes, data)
  
  return data.Data_id, nil
}


// CloseAnnouncement marks an data record to be sell as data not able to be sell anymore
// Admitindo que cada Utilizador guarda na sua base de dados os record_id de todos os anuncios que possui
func (s *SmartContract) CloseAnnouncement(ctx contractapi.TransactionContextInterface, 
	record_id string, owner_id string) error {
	
  dataAsBytes, err := ctx.GetStub().GetState(record_id)

  if err != nil {
      return fmt.Errorf("Failed to read from world state. %s", err.Error())
  }
  
  if dataAsBytes == nil {
	  return fmt.Errorf("%s does not exist", record_id)
  }

  data := new(Data)
  _ = json.Unmarshal(dataAsBytes, data)

  data.Owner_id = owner_id
  dataAsBytes2, _ := json.Marshal(data)

  return ctx.GetStub().PutState(record_id, dataAsBytes2)
}


// ===========================================================================================
// constructQueryResponseFromIterator constructs a array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]Data, error) {
	var results []Data
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		
		data := new(Data)
		_ = json.Unmarshal(queryResponse.Value, data)
		
		results = append(results, *data)
		
	}

	return results, nil
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]Data, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetAnnounces returns the list to the data sell announces in the world state
func (s *SmartContract) GetAnnounces(ctx contractapi.TransactionContextInterface) ([]Data, error) {
   queryString := "{\"selector\":{\"owner\":\"\"}}"
   
   queryResults, err := getQueryResultForQueryString(ctx.GetStub(), queryString)
   if err != nil {
	   return nil, err
   }	
 
   return queryResults, nil
}

// QueryCar returns the Data stored in the world state with given id
func (s *SmartContract) QueryData(ctx contractapi.TransactionContextInterface, record_id string) (*Data, error) {
  dataAsBytes, err := ctx.GetStub().GetState(record_id)

  if err != nil {
      return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
  }

  if dataAsBytes == nil {
      return nil, fmt.Errorf("%s does not exist", record_id)
  }

  data := new(Data)
  _ = json.Unmarshal(dataAsBytes, data)

  return data, nil
}

func main() {

  chaincode, err := contractapi.NewChaincode(new(SmartContract))

  if err != nil {
      fmt.Printf("Error create DataMarket chaincode: %s", err.Error())
      return
  }

  if err := chaincode.Start(); err != nil {
      fmt.Printf("Error starting DataMarket chaincode: %s", err.Error())
  }
}
