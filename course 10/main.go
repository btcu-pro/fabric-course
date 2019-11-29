package main

import (
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "fmt"
    "github.com/hyperledger/fabric/protos/peer"
    "encoding/json"
    "bytes"
)

type CouchDBChaincode struct {

}

type CarStruct struct {
	ObjectType    string    `json:"docType"`
	CarId		  string 	`json:"carid"`
    Owner         string    `json:"owner"`
    Brand         string    `json:"brand"`
    CarName       string    `json:"carname"`
    Price         string    `json:"price"`
}


func (t *CouchDBChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response  {
    return shim.Success(nil)
}


func (t *CouchDBChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response  {
    fun, args := stub.GetFunctionAndParameters()
    if fun == "carInit" {
        return carInit(stub, args)
    } else if fun == "queryCars" {
        return queryCars(stub, args)
    } else if fun == "invokeCars" {
        return invokeCars(stub, args)
    }

    return shim.Error("非法操作, 指定的函数名无效")
}

// 初始化汽车数据
func carInit(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
    car1 := CarStruct{
		ObjectType: "carObj",
		CarId: "car001",
		Owner: "p1",
		Brand: "brand-1",
		CarName: "car-wife",
		Price: "20",
    }

    carByte1, _ := json.Marshal(car1)
    err := stub.PutState(car1.CarId, carByte1)
    if err != nil {
        return shim.Error("初始化第一个汽车失败: "+ err.Error())
	}
	
	car2 := CarStruct{
		ObjectType: "carObj",
		CarId: "car002",
		Owner: "p1",
		Brand: "brand-2",
		CarName: "car-me",
		Price: "40",
    }

    carByte2, _ := json.Marshal(car2)
    err = stub.PutState(car2.CarId, carByte2)
    if err != nil {
        return shim.Error("初始化第二个汽车失败: "+ err.Error())
	}
	
	car3 := CarStruct{
		ObjectType: "carObj",
		CarId: "car003",
		Owner: "p1",
		Brand: "audi-3",
		CarName: "car-son",
		Price: "20",
    }

    carByte3, _ := json.Marshal(car3)
    err = stub.PutState(car3.CarId, carByte3)
    if err != nil {
        return shim.Error("初始化第三个汽车失败: "+ err.Error())
	}
	
	car4 := CarStruct{
		ObjectType: "carObj",
		CarId: "car004",
		Owner: "p2",
		Brand: "brand-1",
		CarName: "car-me",
		Price: "30",
    }

    carByte4, _ := json.Marshal(car4)
    err = stub.PutState(car4.CarId, carByte4)
    if err != nil {
        return shim.Error("初始化第四个汽车失败: "+ err.Error())
    }


    return shim.Success([]byte("初始化汽车成功"))
}



func queryCars(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

    // 查询数据
    result, err := getCarsByQueryString(stub, queryString)
    if err != nil {
        return shim.Error("根据持票人的证件号码批量查询持票人的持有票据列表时发生错误: " + err.Error())
    }
    return shim.Success(result)
}


func getCarsByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

    iterator, err := stub.GetQueryResult(queryString)
    if err != nil {
        return nil, err
    }
    defer  iterator.Close()

    var buffer bytes.Buffer
    var isSplit bool
    for iterator.HasNext() {
        result, err := iterator.Next()
        if err != nil {
            return nil, err
        }

        if isSplit {
            buffer.WriteString("; ")
        }

        buffer.WriteString("key:")
        buffer.WriteString(result.Key)
        buffer.WriteString(", Value: ")
        buffer.WriteString(string(result.Value))

        isSplit = true

    }

    return buffer.Bytes(), nil

}

func invokeCars(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
        return shim.Error("参数数量错误")
    }
	car := CarStruct{
		ObjectType: "carObj",
		CarId: args[0],
		Owner: args[1],
		Brand: args[2],
		CarName: args[3],
		Price: args[4],
    }

    carByte, _ := json.Marshal(car)
    err := stub.PutState(car.CarId, carByte)
    if err != nil {
        return shim.Error("修改汽车失败: "+ err.Error())
	}
	return shim.Success([]byte("修改汽车成功"))

}

func main() {
    err := shim.Start(new(CouchDBChaincode))
    if err != nil {
        fmt.Errorf("启动链码失败: %v", err)
    }
}