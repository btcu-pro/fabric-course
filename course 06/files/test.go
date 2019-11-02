package main

import (
"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/protos/peer"
"fmt"
)

type SimpleChaincode struct {
}


func main(){
    err := shim.Start(new(SimpleChaincode))
    if err != nil{
        fmt.Printf("启动 SimpleChaincode 时发生错误: %s", err)
    }
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{
    args := stub.GetStringArgs()
    if len(args) != 2{
        return shim.Error("初始化的参数只能为2个， 分别代表名称与状态数据")
    }
    err := stub.PutState(args[0], []byte(args[1]))
    if err != nil{
        return shim.Error("在保存状态时出现错误")
    }
    return shim.Success(nil)
}


func (t * SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	fun, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fun == "set"{
		result, err = set(stub, args)
	}else{
		result, err = get(stub, args)
	}
	if err != nil{
		return shim.Error(err.Error())
	}
	
	return shim.Success([]byte(result))
}

func set(stub shim.ChaincodeStubInterface, args []string)(string, error){

    if len(args) != 2{
        return "", fmt.Errorf("给定的参数个数不符合要求")
    }

    err := stub.PutState(args[0], []byte(args[1]))
    if err != nil{
        return "", fmt.Errorf(err.Error())
    }
    return string(args[0]), nil

}


func get(stub shim.ChaincodeStubInterface, args []string)(string, error){
    if len(args) != 1{
        return "", fmt.Errorf("给定的参数个数不符合要求")
    }
    result, err := stub.GetState(args[0])
    if err != nil{
        return "", fmt.Errorf("获取数据发生错误")
    }
    if result == nil{
        return "", fmt.Errorf("根据 %s 没有获取到相应的数据", args[0])
    }
    return string(result), nil

}
