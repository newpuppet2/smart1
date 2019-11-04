package main

/**
 * Shows how to use the history
 **/

import (
	// For printing messages on console
	"fmt"

	// The shim package
	"github.com/hyperledger/fabric/core/chaincode/shim"

	// // peer.Response is in the peer package
	"github.com/hyperledger/fabric/protos/peer"

	// JSON Encoding
	"encoding/json"

	// KV Interface
	// "github.com/hyperledger/fabric/protos/ledger/queryresult"

	"strconv"
)

// LC Represents our chaincode object
type LC struct {
}

// LetterCredit Represents our car asset
type LetterCredit struct {
	DocType			                   string  `json:"docType"`
	Date			                   string  `json:"date"`
	ImporterName 		               string   `json:"importerName"`
	ExporterName		               string   `json:"exporterName"`
	ImporterBankName                   string  `json:"importerBankName"`
	ExporterBankName			       string  `json:"exporterBankName"`
	ProductOrderId			           string    `json:"productOrderId"`
	ProductOrderDetails		           string  `json:"productOrderDetails"`	
	PaymentAmount                     int    `json:"paymentAmount"`
	State                              string  `json:"state"`
	Pendingstate                       string  `json:"pendingstate"`
}

type Letterstatus struct {
	importer			               int  `json:"importer"`
	importerbank		               int   `json:"importerbank"`
	exporterbank		               int   `json:"exporterbank"`
	exporter                           int  `json:"exporter"`
	exportcustoms			           int  `json:"exportcustoms"`
	importcustoms			           int    `json:"importcustoms"`
}

// DocType Represents the object type
const	DocType	= "LC"

const   OwnerPrefix="owner."

func (tradefinance *LC) Init(stub shim.ChaincodeStubInterface) peer.Response {

	// Simply print a message
	counter := 0

	letter1 := Letterstatus{importer: counter, importerbank: counter, exporterbank: counter, exporter: counter, exportcustoms: counter, importcustoms: counter}
	jsonletter, _ := json.Marshal(letter1)
	// Key = VIN#, Value = Car's JSON representation
	stub.PutState("token",jsonletter)
 


	// stub.PutState(counter1, counter)
	fmt.Println("Init executed in tradefinance")

	// Return success
	return shim.Success(nil)
}

func (tradefinance *LC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Get the function name and parameters
	funcName, args := stub.GetFunctionAndParameters()

	if funcName == "CreateLC" {
		// Creates the LC
		return  tradefinance.CreateLC(stub, args)

	} 
	
	else if funcName == "ApproveTrade" {
		// Invoke this function to approve or reject the lc
		return tradefinance.ApproveTrade(stub, args)

	}

/* else if funcName == "Getlc" {
		// Query this function to get txn history for specific vehicle
		return tradefinance.Getlc(stub, args)

	} 
	*/


	// This is not good
	return shim.Error(("Bad Function Name = !!!"))
}

func (tradefinance *LC) CreateLC(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
	bytes,_ := stub.GetState("token")
/*	if err != nil {
		return errorResponse("error", 703)
	}
*/
	// Read the JSON and Initialize the struct
	var ls Letterstatus
	_ = json.Unmarshal(bytes, &ls)
	
	ls.importer = ls.importer + 1
	key := OwnerPrefix+string(ls.importer)
	
	jsonletter, _ := json.Marshal(ls)
	stub.PutState("token", jsonletter)
	
	date := string(args[0])
	importer := string(args[1])
	exporter := string(args[2])
	importerbank := string(args[3])
	exporterbank := string(args[4])
    productdes := string(args[5])
    payment,_ := strconv.Atoi(string(args[6]))
	status := string(args[7])
	pendingstate := string(args[8])
		
	letter := LetterCredit{DocType: DocType, Date: date, ImporterName: importer, ExporterName: exporter, ImporterBankName: importerbank, ExporterBankName: exporterbank, ProductOrderId: key, ProductOrderDetails: productdes, PaymentAmount: payment, State: status,Pendingstate: pendingstate}
	jsonletter, _ = json.Marshal(letter)
	// Key = VIN#, Value = Car's JSON representation
	stub.PutState(key, jsonletter)

    return shim.Success([]byte(jsonletter))
}


func (tradefinance *LC) ApproveTrade(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	firstarg := string(args[0])
        bytes, _ := stub.GetState(firstarg)
	if bytes == nil {
		return shim.Error("Provided ID not found!!!")
	}

	var lc1  LetterCredit
	_ = json.Unmarshal(bytes, &lc1)

	
	app := string(args[1])
    
    

	if app == "reject" {
		lc1.State = "rejected"
		jsonletter, _ := json.Marshal(lc1)
		stub.PutState(firstarg, jsonletter)
		return shim.Success([]byte("Transaction rejected"))
	}
 
	if app == "accept" && lc1.Pendingstate == "importer" {
		lc1.status = "pending"
		lc1.Pendingstate = "importerbank"
		jsonletter, _ := json.Marshal(lc1)
		stub.PutState(firstarg, jsonletter)

		return shim.Success([]byte("Transaction approved"))
	}

	if app == "accept" && lc1.Pendingstate == "importerbank" {
		lc1.status = "pending"
		lc1.Pendingstate = "exporterbank"
		jsonletter, _ := json.Marshal(lc1)
		stub.PutState(args[0], jsonletter)
		return shim.Success([]byte("Transaction approved"))
	}

	if app == "accept" && lc1.Pendingstate == "exporterbank" {
		lc1.status = "pending"
		lc1.Pendingstate = "exporter"
		jsonletter, _ := json.Marshal(lc1)
		stub.PutState(args[0], jsonletter)
		return shim.Success([]byte("Transaction approved"))
	}

	if app == "accept" && lc1.Pendingstate == "exporter" {
		lc1.status = "pending"
		lc1.Pendingstate = "exportcustoms"
		jsonletter, _ := json.Marshal(lc1)
		stub.PutState(args[0], jsonletter)
		return shim.Success([]byte("Transaction approved"))
	}

	if app == "accept" && lc1.Pendingstate == "exportcustoms" {
		lc1.status = "pending"
		lc1.Pendingstate = "importcustoms"
		jsonletter, _ := json.Marshal(lc1)
		stub.PutState(args[0], jsonletter)
		return shim.Success([]byte("Transaction approved"))
	}

	if app == "accept" && lc1.Pendingstate == "importcustoms" {
		lc1.status = "complete"
		lc1.Pendingstate = "importcustoms"
        jsonletter, _ := json.Marshal(lc1)
		stub.PutState(args[0], jsonletter)
		return shim.Success([]byte("Transaction approved"))
	}

}

/*
func (tradefinance *LC) Getlc(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	// Check the number of args
	if len(args) < 1 {
		return shim.Error("MUST provide start & end Key!!")
	}

	bytes := string(args[0])
    bytes1 := bytes + "a" 
	// Get the data by range
	QryIterator, err := stub.GetStateByRange(bytes, bytes1)
	if err != nil {
		fmt.Printf("Error=" + err.Error())
		return shim.Error(err.Error())
	}

	var counterb = 0
	var resultJSON = "["
	// Read the returned results in a loop
	for QryIterator.HasNext() {

		// Hold pointer to the query result
		var resultKV *queryresult.KV
		var err error

		// Get the next element
		resultKV, err = QryIterator.Next()

		if err != nil {
			fmt.Println("Err=" + err.Error())
		} else {
			// Create the data string
			var tokenData string
			tokenData = "{\"key\":\"" + resultKV.GetKey() + "\",\"lc\":" + string(resultKV.GetValue()) + "}"
			if counterb > 0 {
				resultJSON += "," + "\n "
			}
			resultJSON += tokenData
		}
		// Increment counter
		counterb++
	}
	resultJSON += "]"
	resultJSON = "{ \"count\":" + strconv.Itoa(counterb) + ",\"queryResult\":" + resultJSON + "}"

	// Close the query iterator instance
	QryIterator.Close()

	return shim.Success([]byte(resultJSON))
}



func AddData(stub shim.ChaincodeStubInterface,date string,importer string,exporter string,importerbank string,exporterbank string,productid uint,productdes string,payment uint,status string,pendingstate string) {
	letter := LetterCredit{DocType: DocType, Date: date, ImporterName: importer, ExporterName: exporter, ImporterBankName: importerbank, ExporterBankName: exporterbank, ProductOrderId: productid, ProductOrderDetails: productdes, PaymentAmount: payment, State: status,Pendingstate: pendingstate}
	jsonletter, _ := json.Marshal(letter)
	// Key = VIN#, Value = Car's JSON representation
	stub.PutState(productid, jsonletter)
}


*/

// Chaincode registers with the Shim on startup
func main() {
	fmt.Printf("Started Chaincode. LC\n")
	err := shim.Start(new(LC))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}


