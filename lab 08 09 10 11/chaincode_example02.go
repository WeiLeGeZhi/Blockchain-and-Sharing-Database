package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init")
	_, args := stub.GetFunctionAndParameters()
	var A string    // Entities
	var Aval int // Asset holdings
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("A = %d, Aval = %d\n", A, Aval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	var sumid string="sum";
	var sum int=0;
	sum+=Aval;
	err = stub.PutState(sumid, []byte(strconv.Itoa(sum)))
	if err != nil {
		return shim.Error(err.Error())
	}

	var maxid string="max";
	err = stub.PutState(maxid, []byte(A))
	if err != nil {
		return shim.Error(err.Error())
	}

	var minid string="min";
	err = stub.PutState(minid, []byte(A))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "AddBill" {
		return t.AddBill(stub, args)
	} else if function == "Query" {
		return t.Query(stub, args)
	} else if function == "GetSum" {
		return t.GetSum(stub)
	} else if function == "FindMax" {
		return t.FindMax(stub)
	} else if function == "FindMin" {
		return t.FindMin(stub)
	}

	return shim.Error("Invalid invoke function name.")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) AddBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string    // Entities
	var Aval int // Asset holdings
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("A = %d, Aval = %d\n", A, Aval)
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	//add to sum
	var Sumval int
	Sumvalbytes, err := stub.GetState("sum")
	Sumval, _ =strconv.Atoi(string(Sumvalbytes))
	Sumval+=Aval
	err=stub.PutState("sum",[]byte(strconv.Itoa(Sumval)))

	//update max&&min
	var Maxval int
	Maxidbytes, err := stub.GetState("max")
	Maxvalbytes, err := stub.GetState(string(Maxidbytes))
	Maxval, _ =strconv.Atoi(string(Maxvalbytes))
	if Aval > Maxval{
		err =stub.PutState("max",[]byte(A))
	}

	var Minval int
	Minidbytes, err := stub.GetState("min")
	Minvalbytes, err := stub.GetState(string(Minidbytes))
	Minval, _ =strconv.Atoi(string(Minvalbytes))
	if Aval < Minval{
		err =stub.PutState("min",[]byte(A))
	}

	fmt.Println("add succeed")
	jsonResp := "{AddBill succeed:" + A +" "+strconv.Itoa(Aval)+ "}"
	return shim.Success([]byte(jsonResp))
}


// query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{Name:" + A + ",Amount:" + string(Avalbytes) + "}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}

func (t *SimpleChaincode) GetSum(stub shim.ChaincodeStubInterface) pb.Response {

	var Sumval int
	Sumvalbytes, err := stub.GetState("sum")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for sum \"}"
		return shim.Error(jsonResp)
	}
	Sumval, _ =strconv.Atoi(string(Sumvalbytes))
	jsonResp :="{The Sum is "+strconv.Itoa(Sumval)+"}"
	return shim.Success([]byte(jsonResp))
}

func (t *SimpleChaincode) FindMax(stub shim.ChaincodeStubInterface) pb.Response {

	var Maxval int
	Maxidbytes, err := stub.GetState("max")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for max \"}"
		return shim.Error(jsonResp)
	}
	Maxvalbytes, err := stub.GetState(string(Maxidbytes))
	Maxval, _ =strconv.Atoi(string(Maxvalbytes))
	jsonResp :="{The MaxBill is "+string(Maxidbytes)+" "+strconv.Itoa(Maxval)+"}"
	return shim.Success([]byte(jsonResp))
}

func (t *SimpleChaincode) FindMin(stub shim.ChaincodeStubInterface) pb.Response {

	var Minval int
	Minidbytes, err := stub.GetState("min")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for min \"}"
		return shim.Error(jsonResp)
	}
	Minvalbytes, err := stub.GetState(string(Minidbytes))
	Minval, _ =strconv.Atoi(string(Minvalbytes))
	jsonResp :="{The MinBill is "+string(Minidbytes)+" "+strconv.Itoa(Minval)+"}"
	return shim.Success([]byte(jsonResp))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
