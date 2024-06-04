package smallbank

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SmallBankChainCode struct{}

// Init 初始化Smallbank，这里初始化两个账户的savings和checkings
// 所有账户的savings用"账户_savings"表示，checkings用"账户_checkings"表示
// 初始化需要6个参数, 账户1 savings余额 checkings余额 账户2 savings余额 checkings余额
func (t *SmallBankChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// Initialize the chaincode
	accountA := args[0]
	accountASavings, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	accountACheckings, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// 将账户A的相关信息存入
	// 存入savings
	err = stub.PutState(accountA+"_savings", []byte(strconv.Itoa(accountASavings)))
	// 存入checkings
	err = stub.PutState(accountA+"_checkings", []byte(strconv.Itoa(accountACheckings)))

	accountB := args[3]
	accountBSavings, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	accountBCheckings, err := strconv.Atoi(args[5])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// 存入savings
	err = stub.PutState(accountB+"_savings", []byte(strconv.Itoa(accountBSavings)))
	if err != nil {
		return shim.Error(err.Error())
	}
	// 存入checkings
	err = stub.PutState(accountB+"_checkings", []byte(strconv.Itoa(accountBCheckings)))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("Init Account ", accountA, " ", accountB)
	return shim.Success(nil)
}
func (t *SmallBankChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "create" {
		// 创建一个新的账户
		return t.create(stub, args)
	} else if function == "transactSavings" {
		// 向储蓄账户增加一定余额
		return t.transactSavings(stub, args)
	} else if function == "depositChecking" {
		// 向支票账户增加一定余额
		return t.depositChecking(stub, args)
	} else if function == "sendPayment" {
		// 在两个支票账户间转账
		return t.sendPayment(stub, args)
	} else if function == "writeCheck" {
		// 减少一个支票账户
		return t.writeCheck(stub, args)
	} else if function == "amalgamate" {
		// 将储蓄账户的资金全部转到支票账户
		return t.amalgamate(stub, args)
	} else if function == "query" {
		// 读取一个用户的支票账户以及储蓄账户
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name.")
}

// create 创建一个账户，根据传入的参数初始化saving、checking
// args: [账户, saving, checking]
func (t *SmallBankChainCode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	account := args[0]

	// todo
	savings, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	checkings, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	err = stub.PutState(account+"_savings", []byte(strconv.Itoa(savings)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(account+"_checkings", []byte(strconv.Itoa(checkings)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Create Account ", account)
	return shim.Success(nil)
}

// 向储蓄账户增加一定余额
// 参数: [账户，金额]， 这里只需要账户, saving由"账户_savings"拼接得到
func (t *SmallBankChainCode) transactSavings(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	account := args[0]

	// todo
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for amount")
	}

	savingsBytes, err := stub.GetState(account + "_savings")
	if err != nil {
		return shim.Error("Failed to get savings state")
	}
	if savingsBytes == nil {
		return shim.Error("Savings account not found")
	}
	savings, _ := strconv.Atoi(string(savingsBytes))
	savings += amount

	err = stub.PutState(account+"_savings", []byte(strconv.Itoa(savings)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Update Account ", account)
	return shim.Success(nil)
}

// depositChecking 向支票账户增加一定余额
// 参数: [账户，金额]， 这里只需要账户, checking由"账户_checkings"拼接得到
func (t *SmallBankChainCode) depositChecking(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	account := args[0]

	// todo
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for amount")
	}

	checkingsBytes, err := stub.GetState(account + "_checkings")
	if err != nil {
		return shim.Error("Failed to get checkings state")
	}
	if checkingsBytes == nil {
		return shim.Error("Checkings account not found")
	}
	checkings, _ := strconv.Atoi(string(checkingsBytes))
	checkings += amount

	err = stub.PutState(account+"_checkings", []byte(strconv.Itoa(checkings)))
	if err != nil {
		return shim.Error(err.Error())
	}


	fmt.Println("Update Account ", account)
	return shim.Success(nil)
}

// sendPayment 在两个支票账户间转账
// 参数: [accountA, accountB, amount], A向B转账amount(在checkings里)
func (t *SmallBankChainCode) sendPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	accountA, accountB := args[0], args[1]

	// todo
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for amount")
	}

	checkingsABytes, err := stub.GetState(accountA + "_checkings")
	if err != nil {
		return shim.Error("Failed to get checkings state")
	}
	if checkingsABytes == nil {
		return shim.Error("Checkings account A not found")
	}
	checkingsA, _ := strconv.Atoi(string(checkingsABytes))

	checkingsBBytes, err := stub.GetState(accountB + "_checkings")
	if err != nil {
		return shim.Error("Failed to get checkings state")
	}
	if checkingsBBytes == nil {
		return shim.Error("Checkings account B not found")
	}
	checkingsB, _ := strconv.Atoi(string(checkingsBBytes))

	if checkingsA < amount {
		return shim.Error("Insufficient funds in checkings account A")
	}

	checkingsA -= amount
	checkingsB += amount

	err = stub.PutState(accountA+"_checkings", []byte(strconv.Itoa(checkingsA)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(accountB+"_checkings", []byte(strconv.Itoa(checkingsB)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Update Account ", accountA, " ", accountB)
	return shim.Success(nil)
}

// 减少一个支票账户
// 参数: [账户，金额]，扣除一笔钱
func (t *SmallBankChainCode) writeCheck(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	account := args[0]
	// todo
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for amount")
	}

	checkingsBytes, err := stub.GetState(account + "_checkings")
	if err != nil {
		return shim.Error("Failed to get checkings state")
	}
	if checkingsBytes == nil {
		return shim.Error("Checkings account not found")
	}
	checkings, _ := strconv.Atoi(string(checkingsBytes))

	if checkings < amount {
		return shim.Error("Insufficient funds in checkings account")
	}

	checkings -= amount

	err = stub.PutState(account+"_checkings", []byte(strconv.Itoa(checkings)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Update Account ", account)
	return shim.Success(nil)
}

// amalgamate 将储蓄账户的资金全部转到支票账户
// 参数: [account]
func (t *SmallBankChainCode) amalgamate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	account := args[0]
	// todo
	savingsBytes, err := stub.GetState(account + "_savings")
	if err != nil {
		return shim.Error("Failed to get savings state")
	}
	if savingsBytes == nil {
		return shim.Error("Savings account not found")
	}
	savings, _ := strconv.Atoi(string(savingsBytes))

	checkingsBytes, err := stub.GetState(account + "_checkings")
	if err != nil {
		return shim.Error("Failed to get checkings state")
	}
	if checkingsBytes == nil {
		return shim.Error("Checkings account not found")
	}
	checkings, _ := strconv.Atoi(string(checkingsBytes))

	checkings += savings
	savings = 0

	err = stub.PutState(account+"_checkings", []byte(strconv.Itoa(checkings)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(account+"_savings", []byte(strconv.Itoa(savings)))
	if err != nil {
		return shim.Error(err.Error())
	}


	fmt.Println("Update Account ", account)
	return shim.Success(nil)
}

// query 查询一个账户的余额信息
// 参数: [account]
func (t *SmallBankChainCode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	account := args[0]
	savings, err := stub.GetState(account + "_savings")
	if err != nil {
		return shim.Error("Failed to get savings state")
	}
	checkings, err := stub.GetState(account + "_checkings")
	if err != nil {
		shim.Error("Failed to get checkings state")
	}
	//fmt.Println("Account ", account)
	//fmt.Println("	Savings: ", savingsValue)
	//fmt.Println("	Checkings: ", checkingsValue)
	// shim.LogInfo("This is an info log message from the chaincode")
	return shim.Success([]byte("Saving: " + string(savings) + ", Checking: " + string(checkings)))
}
