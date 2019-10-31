package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

func main()  {
	err := shim.Start(new(HelloChaincode))
	if err != nil {
		fmt.Printf("链码启动失败: %v", err)
	}
}

type HelloChaincode struct {

}

func (t *HelloChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response  {
	fmt.Println("开始实例化链码....")

	// 获取参数
	//args := stub.GetStringArgs()
	_, args := stub.GetFunctionAndParameters()
	// 判断参数长度是否为2个
	if len(args) != 2 {
		return shim.Error("指定了错误的参数个数")
	}

	fmt.Println("保存数据......")

	// 通过调用PutState方法将数据保存在账本中
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error("保存数据时发生错误...")
	}

	fmt.Println("实例化链码成功")

	return shim.Success(nil)

}

func (t *HelloChaincode)  Invoke(stub shim.ChaincodeStubInterface) peer.Response  {
    // 获取调用链码时传递的参数内容(包括要调用的函数名及参数)
    fun, args := stub.GetFunctionAndParameters()

    // 客户意图
    if fun == "query"{
        return query(stub, args)
    }

        return shim.Error("非法操作, 指定功能不能实现")
}

func query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 检查传递的参数个数是否为1
	if len(args) != 1{
		return shim.Error("指定的参数错误，必须且只能指定相应的Key")
	}

	// 根据指定的Key调用GetState方法查询数据
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据指定的 " + args[0] + " 查询数据时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的 " + args[0] + " 没有查询到相应的数据")
	}

	// 返回查询结果
	return shim.Success(result)
}
