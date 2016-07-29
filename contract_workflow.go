package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const WORKFLOW_TEMPLATE_KEY string = "WORKFLOW_TEMPLATE_KEY"

const MYVERSION string = "1.0.0"

type ContractState struct {
	Version            string `json:"version"`
	WorkflowTemplateCC string `json:"workflowTemplateCC"`
}

//Event represents events that occur during the workflow of a document
type Event struct {
	User         string `json:"user"`
	Role         string `json:"role"`
	Organization string `json:"organization"`
	DocId        string `json:docId`
	DocType      string `json:docType`
	TimeStamp    string `json:"timeStamp"`
	EventType    string `json:"eventType"`
	DocHash      string `json:"docHash"` //do I need this for all events like read?
	DocUrl       string `json:docUrl`
}

//The following struct definitions should match Sayandeep's struct definition
type ActionId struct {
	Id string `json:"id"`
}
//Specification for all documents needs to be in this format
type DocType struct {
	Id 						string 	`json: "id"`
	Name 						string `json: "name"`
	Type 						string `json: "type"`
	Version 					string `json: "version"`
	AllowedActTypeList 				[]ActionId `json: "allowedActTypeList"`
}
//This is a role think exporter, importer etc. who can take Actions
type Actor struct {
	Id 		string `json: "Id"`
	Name 	string `json: "name"`
}
//Action is a combination of document and a allowable action on the document
//Specification for all actions needs to be in this format
type ActType struct {
	Id 					string `json: "id"`
	Name 			string `json: "name"`
	Type string `json: "type"`
}

//need checks to ensure that the documents and actions are valid
type Action struct {
	Id 					string 	`json: "id"`
	DocTypeId 			string 	`json: "docTypeId"`//Name of the DocType on which action is defined
	ActTypeId 			string 	`json: "actTypeId"`//Name of the ActType which is defined on the document
	ActorId				string 	`json: "actorId"`	//Name of the Actor who is supposed to undertake the action
}

type DependencyInfo struct {
	Id 					string `json: "id"`
	DependencyList		[]ActionId `json:"dependencyInfo"`
}

type WorkflowIO struct {
	Id             string           `json: "id"`
	Name           string           `json: "name"`
	Version        string           `json: "version"`
	Desc           string           `json: "desc"`
	ActTypeList    []ActType        `json: "actTypeList"`    //List of ActionTypes we will need in this workflow
	DocTypeList    []DocType        `json: "docTypeList"`    //List of DocumentTypes we will need in this workflow.
	ActionList     []Action         `json: "actionList"`     //List of Actions which will be performed in this workflow
	ActorList      []Actor          `json: "actorList"`      //List of Actors in this workflow
	AclList        []Action         `json: "aclList"`        //List of ACL information
	DependencyList []DependencyInfo `json: "dependencyList"` //array of actiondependencyInfo
}

//Sayandeep struct defintions end

type SimpleChaincode struct {
}

type RoleToInstanceMapping struct {
	InstanceMapping map[string]string
}

type Contract struct {
	WorkflowId                       string             `json:"WorkflowId"`
	ContractName                     string             `json:"contractName"`
	ContractId                       string             `json:contractId` //contractId
	Documents                        []string           `json:documents`  //Array of doc ids
	Events                           map[string][]Event `json:"events"`
	ParticipantRoleToInstanceMapping map[string]string  `json:"participants"`
	DocTypeToDocIdMapping            map[string]string  `json:docTypeToDocIdMapping` //there can be only one document of each type?
	//All participants should be in a map to make it extensible
	/*
		Exporter string `json:"exporter"`
		Importer string `json:"importer"`
		FreightForwarder string `json:"freighForwarder"`
		ExportRevenue string `json:"exportRevenue"`
		ExportPort string `json:"exportPort"`
		ExportCustoms string `json:"exportCustoms"`
		ImportCustoms string `json:"importCustoms"`
	*/
}

/*
	Chaincode initialization code.
*/
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//Initialize chaincode--may be use a version number
	var stateArg ContractState
	var err error

	if len(args) != 1 {
		return nil, errors.New("init expects one argument, a JSON string with tagged version string and chaincode uuid for the workflow template chaincode")
	}
	err = json.Unmarshal([]byte(args[0]), &stateArg)
	if err != nil {
		return nil, errors.New("Version argument unmarshal failed: " + fmt.Sprint(err))
	}
	if stateArg.Version != MYVERSION {
		return nil, errors.New("Contract version " + MYVERSION + " must match version argument: " + stateArg.Version)
	}
	// set the chaincode uuid of the compliance contract
	// to the global variable

	if stateArg.WorkflowTemplateCC == "" {
		return nil, errors.New("Workflow Template chaincode id is mandatory")
	}

	contractStateJSON, err := json.Marshal(stateArg)
	if err != nil {
		return nil, errors.New("Marshal failed for contract state" + fmt.Sprint(err))
	}
	err = stub.PutState(WORKFLOW_TEMPLATE_KEY, contractStateJSON)
	if err != nil {
		return nil, errors.New("Contract state failed PUT to ledger: " + fmt.Sprint(err))
	}

	return nil, nil
}

/*
	Called during invoke calls of chaincode
*/

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	fmt.Println("Invoke: Function called is %s", function)

	if function == "addContract" {
		t.AddContract(stub, args)
	} else if function == "addSignature" {
		t.AddSignature(stub, args)
	} else if function == "issue" {
		t.IssueDocument(stub, args)
	} else if function == "read" {
		t.ReadDocument(stub, args)
	} else if function == "addCert" {
		t.AddCert(stub, args)
	} else if function == "addSampleWfIOJson" {
		t.AddSampleWfIOJson(stub,args)
	}else if function == "printSampleWfIOJson" {
		t.PrintSampleWfIOJson(stub,args)
	}
	return nil, nil

}

//addCert	adds a certificate to the access control list of the chaincode
func (t *SimpleChaincode) AddCert(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// args=[what should be the arg strcuture be like]
	return nil, nil
}

/*
	Registers a new contract in the blockchain
	args=["workflowId","contractId", "contractName", {"role":"instance-mapping"} ]
	Inputs:
		workflowId : identifies the workflow type for the given contract
	 	contractId : id of the contract generated in the App server
		contractName : name of the contract entered by the Exporter
		role to instance mapping json : Specifies mappings for a given contractId between each role mentioned in the workflow for given workflowId
																		to specified participants.


*/

func (t *SimpleChaincode) AddContract(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
	//var contract Contract
	workflowId := args[0]
	contractId := args[1]
	contractName := args[2]
	roleToInstanceMappingJson := args[3]
	participants := make(map[string]string)

	//TODO: should check if all participants specified in the workflow are in the contract

	//fmt.Printf("Inputs args are--workflowid: %s ",workflowId," contractId: %s", contractId, "contractName: %s ", contractName )

	fmt.Printf("Inputs args are--workflowid: %s contractId: %s contractName: %s instanceMapping: %s \n", workflowId, contractId, contractName, participants)

	contractBytes, err := stub.GetState(contractId)
	if err != nil {
		return nil, errors.New("Error in getting State from Ledger")
	}
	if contractBytes != nil {
		return nil, errors.New("ContractId already exists")

	}
	var f interface{}
	err = json.Unmarshal([]byte(roleToInstanceMappingJson), &f)
	if err != nil {
		fmt.Println("Error in unmarshaling roleToInstanceMappingJson")
	}
	mapping := f.(map[string]interface{})
	for k, v := range mapping {
		v1 := v.(string)
		participants[k] = v1
	}

	contract := new(Contract)
	contract.WorkflowId = workflowId
	contract.ContractId = contractId
	contract.ContractName = contractName
	contract.ParticipantRoleToInstanceMapping = participants

	contract.Events = make(map[string][]Event)
	contract.DocTypeToDocIdMapping = make(map[string]string)

	_, errValid := t.IsValidContract(stub,contract)
	if errValid != nil {
		return nil, errValid
	}

	fmt.Println("Going to store contractState")
	contractJson, err := json.Marshal(contract)

	err = stub.PutState(contractId, []byte(contractJson))
	if err != nil {
		fmt.Println("Failed PUT to ledger while adding Contract")
		return nil, errors.New("Failed PUT to ledger while adding Contract: " + fmt.Sprint(err))
	}
	return nil, nil
}

//IsValidContract checks the validity of an input contract by verifying with the corresponding workflow specification
func (t *SimpleChaincode) IsValidContract(stub *shim.ChaincodeStub,contract * Contract) (bool, error) {

	fmt.Printf("Inside IsValidContract\n")
	var wfIO WorkflowIO
	workflowId := contract.WorkflowId
	fmt.Printf("workflowId: %s",workflowId)
  roleToInstanceMapping := contract.ParticipantRoleToInstanceMapping

	//TODO: uncomment this
	//	wfIOBytes, err := t.GetWorkFlowTemplate(stub,workflowId)

	wfIOBytes, err := t.GetSampleWfIOJson(stub)
	fmt.Println("Got Sample wfIO as ",string(wfIOBytes))
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(wfIOBytes, &wfIO)
	if err != nil {
		return false, err
	}
	actorList := wfIO.ActorList

	fmt.Println("Entering for loop")
	//check if all roles have instance mapping
	for _,actor := range actorList {
		roleName := actor.Name
		fmt.Printf("Actor rolename is %s",roleName)
		_, ok := roleToInstanceMapping[roleName]
		if ok {
			fmt.Println(roleName," exists in mapping")
		} else {
			fmt.Println(roleName," does not exist in mapping")
			return false, fmt.Errorf("%s is missing in the mapping",roleName)
		}
	}
	fmt.Println("Exiting for loop")

	return true, nil
}

//AddSignature checks if the user/role is authorized and records a signature event to the document by the user
//TODO: ACL check
func (t *SimpleChaincode) AddSignature(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// args=["contractId", "docId", {"User":"Masai","TimeStamp":"18:00 Aug 13, 2016","EventType":"Sign"},"docHash" ]
	return t.RegisterDocumentEvent(stub, args)
	//return nil, nil
}

func (t *SimpleChaincode) IssueDocument(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// args=["contractId", "docId", {"User":"Masai","TimeStamp":"18:00 Aug 13, 2016","EventType":"Sign"},"docHash" ]
	return t.RegisterDocumentEvent(stub, args)
	//return nil, nil
}

//IssueDocument checks if the user is authorized and records the issuance of a document by the user
func (t *SimpleChaincode) RegisterDocumentEvent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// args=["contractId", "docId", {"User":"Masai","Role":"Exporter","TimeStamp":"18:00 Aug 13, 2016","EventType":"Issue",DocId:"doc123","DocType":"Phyto","DocUrl":"bit.ly/abcd","DocHash":"AXSEDFF"}]

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	var contract Contract

	contractId := args[0]
	docId := args[1]
	eventJson := args[2]
	fmt.Println(eventJson)

	var docEvent Event
	err := json.Unmarshal([]byte(eventJson), &docEvent)
	if err != nil {
		return nil, errors.New("document event Json: Unmarshaling error")
	}
	docType := docEvent.DocType
	fmt.Println("Doc type: ", docType)
	eventType := docEvent.EventType
	fmt.Println("Event type: ", eventType)

	eDocId := docEvent.DocId
	if eDocId != docId {
		fmt.Println("docId: DocIds in eventJson and args array dont match")
		err = errors.New("docId: DocIds in eventJson and args array dont match")
		return nil, err
	}

	//TODO:Call to ACL Check

	contractBytes, err := stub.GetState(contractId)
	if err != nil {
		return nil, errors.New("Error in getting State from Ledger")
	}
	if contractBytes == nil {
		return nil, errors.New("ContractId: Could not get state from ledger")
	}
	err = json.Unmarshal(contractBytes, &contract)
	if err != nil {
		return nil, errors.New("Contract: Unmarshaling error")
	}

	if eventType == "Issue" {
		contract.Documents = append(contract.Documents, docId)
		contract.DocTypeToDocIdMapping[docType] = docId
	}

	contract.Events[docId] = append(contract.Events[docId], docEvent)
	contractJson, err := json.Marshal(contract)

	err = stub.PutState(contractId, []byte(contractJson))
	if err != nil {
		fmt.Println("Failed PUT to ledger while adding Contract")
		return nil, errors.New("Failed PUT to ledger while adding Contract: " + fmt.Sprint(err))
	}

	return nil, nil
}

func (t *SimpleChaincode) ReadDocument(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// args=["contractId", "docId", {"User":"Kephis","Role":"ExportAgriculture","TimeStamp":"18:00 Aug 13, 2016","EventType":"Read"}]
	return t.RegisterDocumentEvent(stub, args)
	//return nil, nil
}


func (t *SimpleChaincode) GetEventsForDocument(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("args: Incorrect number of arguments.Expecting 2")
	}
	var contract Contract
	contractId := args[0]
	docId := args[1]
	contractBytes, err := stub.GetState(contractId)
	if err != nil {
		return nil, errors.New("Error in getting State from Ledger")
	}

	if contractBytes != nil {
		err = json.Unmarshal(contractBytes, &contract)
		if err != nil {
			return nil, errors.New("getEventsForDocument: Unmarshaling error of ContractBytes")
		} else {
			fmt.Printf("Contract is %s", contract)
			events := contract.Events[docId]
			if events == nil {
				return nil, errors.New("docId: Does not exist for contract")
			}
			eventsBytes, err1 := json.Marshal(events)
			if err1 != nil {
				return nil, errors.New("eventBytes: Marshaling error")
			} else {
				return eventsBytes, nil
			}
		}
	} else {
		return nil, errors.New("Error in getting ContractBytes from Ledger")
	}

}


//GetWorkFlowTemplate returns the workflow template for a given workflowId
func (t *SimpleChaincode) GetWorkFlowTemplate(stub *shim.ChaincodeStub, workflowId string) ([]byte, error) {
	//workflowId := args[0]
	//TODO: call workflowTemplate chaincode to get template by workflowId.
	//TODO: should return a marshalled WorkFlowIO
	//TODO: wfIOBytes := querychaincode
	var contractState ContractState
	var wfIO WorkflowIO
	var wfIOBytes []byte
	contractStateJSON, err := stub.GetState(WORKFLOW_TEMPLATE_KEY)
	if err != nil {
		return nil, errors.New("Unable to fetch container and compliance contract keys")
	}
	err = json.Unmarshal(contractStateJSON, &contractState)
	if err != nil {
		return nil, err
	}
	workflowTemplateChaincode := contractState.WorkflowTemplateCC

	//TODO: check with Sayandeep on what function to query
	function := "query" //
	args := []string{workflowId}
	wfIOBytes,err = stub.QueryChaincode(workflowTemplateChaincode, function, args)
	if err != nil {
		return nil, errors.New("workflowTemplate: Error in querying from workflow template chaincode")
	} else {
		err = json.Unmarshal(wfIOBytes, &wfIO)
		if err != nil {
			return nil, err
		} else {
			return wfIOBytes, nil
		}
	}
}

//Query chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if function == "getEventsForDocument" {
		return t.GetEventsForDocument(stub, args)
	}
	contractId := args[0]
	var contract Contract
	contractBytes, err := stub.GetState(contractId)
	if err != nil {
		return nil, errors.New("Error in getting State from Ledger")
	}

	if contractBytes != nil {
		err = json.Unmarshal(contractBytes, &contract)
		fmt.Printf("Contract is %s", contract)
		if err != nil {
			return nil, errors.New("Unmarshaling error of ContractBytes")
		} else {
			return contractBytes, nil
		}
	} else {
		return nil, errors.New("Error in getting ContractBytes from Ledger")
	}

	return nil, nil
}



/*

	Testing code

*/

func (t *SimpleChaincode) AddSampleWfIOJson(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	err := stub.PutState(WORKFLOW_TEMPLATE_KEY, []byte(args[0]))
	if err != nil {
		return nil, errors.New("SampleWfIOJson state failed PUT to ledger: " + fmt.Sprint(err))
	}

	return nil, nil
}

func (t *SimpleChaincode) GetSampleWfIOJson(stub *shim.ChaincodeStub) ([]byte, error) {
	WfIOJsonBytes, err := stub.GetState(WORKFLOW_TEMPLATE_KEY)
	if err != nil {
		return nil, errors.New("Unable to fetch WfIOJsonBytes")
	}
	return WfIOJsonBytes,nil
}

func (t *SimpleChaincode) PrintSampleWfIOJson(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	WfIOJsonBytes, err := stub.GetState(WORKFLOW_TEMPLATE_KEY)
	if err != nil {
		return nil, errors.New("Unable to fetch WfIOJsonBytes")
	}
	fmt.Println(string(WfIOJsonBytes))
	return nil,nil
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

/*

func main() {
	fmt.Println("Hello, playground")
	e1 := Event{
		User: "Masai",
		TimeStamp: "AUG 13 2016",
		EventType: "Issue",
	}
	e2 := Event{
		User: "Kephis",
		TimeStamp: "AUG 15 2016",
		EventType: "Sign",
	}
	eArr := []Event{e1,e2}
	eL := make(map[string][]Event)
	eL["doc123"] = eArr
	m := Contract{
	ContractName: "ABC",
	ContractId: "123",
	Events: eL,
	}
	fmt.Println(m)

	b,_ := json.Marshal(m)
	s := string(b)
	fmt.Println(s)
}

*/
