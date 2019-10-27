package main

// 引入相应的依赖包
import (
"fmt"
"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {

}

// 链码实例化（instantiate）或 升级（upgrade）时调用 Init 方法
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{

return shim.Success(nil)
}

// 链码收到调用（invoke） 或 查询 （query）时调用 Invoke 方法
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
return shim.Success(nil)
}

// 主函数 ，调用 shim.Start 方法
func main() {
err := shim.Start(new(SimpleChaincode))

if( err!= nil){
fmt.Printf("Error starting Simple Chaincode is %s \n",err)
}
}
