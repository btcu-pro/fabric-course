# course 06 链码开发与实践
## `Contents`
- [course 06 链码开发与实践](#course-06-%e9%93%be%e7%a0%81%e5%bc%80%e5%8f%91%e4%b8%8e%e5%ae%9e%e8%b7%b5)
  - [`Contents`](#contents)
  - [6.1 如何利用Fabric提供的接口编写链码](#61-%e5%a6%82%e4%bd%95%e5%88%a9%e7%94%a8fabric%e6%8f%90%e4%be%9b%e7%9a%84%e6%8e%a5%e5%8f%a3%e7%bc%96%e5%86%99%e9%93%be%e7%a0%81)
    - [6.1.1 链码接口](#611-%e9%93%be%e7%a0%81%e6%8e%a5%e5%8f%a3)
      - [Init 与 Invoke 方法](#init-%e4%b8%8e-invoke-%e6%96%b9%e6%b3%95)
    - [6.1.2 必要结构](#612-%e5%bf%85%e8%a6%81%e7%bb%93%e6%9e%84)
      - [依赖包](#%e4%be%9d%e8%b5%96%e5%8c%85)
  - [6.2 如何操作账本数据：熟悉链码相关API](#62-%e5%a6%82%e4%bd%95%e6%93%8d%e4%bd%9c%e8%b4%a6%e6%9c%ac%e6%95%b0%e6%8d%ae%e7%86%9f%e6%82%89%e9%93%be%e7%a0%81%e7%9b%b8%e5%85%b3api)
    - [6.2.1 参数解析相关API](#621-%e5%8f%82%e6%95%b0%e8%a7%a3%e6%9e%90%e7%9b%b8%e5%85%b3api)
    - [6.2.2 账本数据状态操作API](#622-%e8%b4%a6%e6%9c%ac%e6%95%b0%e6%8d%ae%e7%8a%b6%e6%80%81%e6%93%8d%e4%bd%9capi)
    - [6.2.3 交易信息相关API](#623-%e4%ba%a4%e6%98%93%e4%bf%a1%e6%81%af%e7%9b%b8%e5%85%b3api)
    - [6.2.4 事件处理API](#624-%e4%ba%8b%e4%bb%b6%e5%a4%84%e7%90%86api)
    - [6.2.5 对 PrivateData 操作的 API](#625-%e5%af%b9-privatedata-%e6%93%8d%e4%bd%9c%e7%9a%84-api)
  - [6.3 链码实现的Hello World](#63-%e9%93%be%e7%a0%81%e5%ae%9e%e7%8e%b0%e7%9a%84hello-world)
    - [6.3.1 链码开发](#631-%e9%93%be%e7%a0%81%e5%bc%80%e5%8f%91)
    - [6.3.2 链码测试](#632-%e9%93%be%e7%a0%81%e6%b5%8b%e8%af%95)
  - [6.4 动手编码一：链码实现资产管理](#64-%e5%8a%a8%e6%89%8b%e7%bc%96%e7%a0%81%e4%b8%80%e9%93%be%e7%a0%81%e5%ae%9e%e7%8e%b0%e8%b5%84%e4%ba%a7%e7%ae%a1%e7%90%86)
    - [6.4.1 资产链码开发](#641-%e8%b5%84%e4%ba%a7%e9%93%be%e7%a0%81%e5%bc%80%e5%8f%91)
    - [6.4.2 链码测试](#642-%e9%93%be%e7%a0%81%e6%b5%8b%e8%af%95)
  - [6.5 动手编码二：链码实现转账](#65-%e5%8a%a8%e6%89%8b%e7%bc%96%e7%a0%81%e4%ba%8c%e9%93%be%e7%a0%81%e5%ae%9e%e7%8e%b0%e8%bd%ac%e8%b4%a6)
    - [6.5.1 转账链码开发](#651-%e8%bd%ac%e8%b4%a6%e9%93%be%e7%a0%81%e5%bc%80%e5%8f%91)
    - [6.5.2 链码测试](#652-%e9%93%be%e7%a0%81%e6%b5%8b%e8%af%95)



## 6.1 如何利用Fabric提供的接口编写链码

开发链码，离不开 Hyperledger Fabric 提供的 SDK(Software Development Kit)，为了方便诸多不同的应用场景且使用不同语言的开发人员，Hyperledger Fabric 提供了许多不同的 SDK 来支持各种编程语言。如：

* Hyperledger Fabric Node SDK：https://github.com/hyperledger/fabric-sdk-node
* Hyperledger Fabric Java SDK：https://github.com/hyperledger/fabric-sdk-java
* Hyperledger Fabric Python SDK：https://github.com/hyperledger/fabric-sdk-py
* Hyperledger Fabric Go SDK：https://github.com/hyperledger/fabric-sdk-go


这节课中我们将使用 Golang 进行链码的开发，所以我们应该确定在本系统中有 Hyperledger Fabric 提供的相关API，其它语言的 SDK 我们不在本课程中进行讨论。

如果本地系统中没有相关的API，请执行如下下载命令：
```shell
go get -u github.com/hyperledger/fabric/core/chaincode/shim
```
如果下载不下来，可以试试
```shell
gopm get -g github.com/hyperledger/fabric/core/chaincode/shim
```

不过仍有同学的 gopm 半天也下载不了，就直接用 git 吧：
```
cd $GOPATH/src/github.com/hyperledger/

git clone https://github.com/hyperledger/fabric.git
```

**可能需要等一阵子**

### 6.1.1 链码接口
链码启动必须通过调用 shim 包中的 Start 函数，而 Start 函数被调用时需要传递一个类型为 Chaincode 的参数，这个参数 Chaincode 是一个接口类型，该接口中有两个重要的函数 Init 与 Invoke 。

Chaincode 接口定义如下：
```go
type Chaincode interface{
    Init(stub ChaincodeStubInterface) peer.Response
    Invoke(stub ChaincodeStubInterface) peer.Response
}
```
#### Init 与 Invoke 方法

编写链码，关键是实现 Init 与 Invoke 两个方法，必须由所有链码实现。Fabric 通过调用指定的函数来运行事务。

* **Init**：在链码实例化或升级时被调用, 完成初始化数据的工作。
* **invoke**：更新或查询提案事务中的分类帐本数据状态时，Invoke 方法被调用， 因此响应调用或查询的业务实现逻辑都需要在此方法中编写实现。  
  
在实际开发中，开发人员可以自行定义一个结构体，然后重写 Chaincode 接口的两个方法，并将两个方法指定为自定义结构体的成员方法；具体可参考下面 6.2 的内容。

### 6.1.2 必要结构
#### 依赖包

shim 包为链码提供了 API 用来访问/操作数据状态、事务上下文和调用其他链代码；peer 包提供了链码执行后的响应信息。所以开发链码需要引入如下依赖包：

* "github.com/hyperledger/fabric/core/chaincode/shim"
  * shim 包提供了链码与账本交互的中间层。
  * 链码通过 shim.ChaincodeStub 提供的方法来读取和修改账本的状态。
* "github.com/hyperledger/fabric/protos/peer"
  * peer.Response：封装的响应信息。

一个开发的链码源文件的必要结构如下：
```golang
package main

// 引入必要的包
import(
    "fmt"   // 标准库的格式化输入输出包

    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)

// 声明一个结构体
type SimpleChaincode struct {

}

// 为结构体添加Init方法
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{
  // 在该方法中实现链码初始化或升级时的处理逻辑
  // 编写时可灵活使用stub中的API
}

// 为结构体添加Invoke方法
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
  // 在该方法中实现链码运行中被调用或查询时的处理逻辑
  // 编写时可灵活使用stub中的API
}

// 主函数，需要调用shim.Start（ ）方法
func main() {
  err := shim.Start(new(SimpleChaincode))
  if err != nil {
     fmt.Printf("Error starting Simple chaincode: %s", err)
  }
}
```

> 因为链码是一个可独立运行的应用，所以必须声明在一个 main 包中，并且提供相应的 main 函数做为应用入口。


## 6.2 如何操作账本数据：熟悉链码相关API
> 现在我们知道了编写链码的基本接口及所需要的结构，那么实际中对账本数据该如何在什么情况下调用什么 API 进行操作？

shim 包提供给链码的相应接口有如下几种类型：

* **参数解析 API**：调用链码时需要给被调用的目标函数/方法传递参数，与参数解析相关的 API 提供了获取这些参数（包含被调用的目标函数/方法名称）的方法。
* **账本状态数据操作 API**：该类型的 API 提供了对账本数据状态进行操作的方法，包括对状态数据的查询及事务处理等。
* **交易信息获取 API**：获取提交的交易信息的相关 API。
* **事件处理 API**：与事件处理相关的 API。
* **对 PrivateData 操作的 API**： Hyperledger Fabric 在 1.2.0 版本中新增的对私有数据操作的相关 API。
下面我们介绍每一种类型相对应的 API 的定义及调用时所需参数。

### 6.2.1 参数解析相关API
`GetArgs() [][]byte`：返回调用链码时在交易提案中指定提供的被调用函数及参数列表

`GetArgsSlice() ([]byte, error)`：返回调用链码时在交易提案中指定提供的参数列表

`GetFunctionAndParameters() (function string, params []string)`：返回调用链码时在交易提案中指定提供的被调用的函数名称及其参数列表

`GetStringArgs() []string`：返回调用链码时指定提供的参数列表

> 在实际开发中，常用的获取被调用函数及参数列表的API一般为： GetFunctionAndParameters() 及 GetStringArgs() 两个。

### 6.2.2 账本数据状态操作API
`GetState(key string) ([]byte, error)` ：根据指定的 Key 查询相应的数据状态。

`PutState(key string, value []byte) error`：根据指定的 key，将对应的 value 保存在分类账本中。

`DelState(key string) error`：根据指定的 key 将对应的数据状态删除

`GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)`：根据指定的开始及结束 key，查询范围内的所有数据状态。注意：结束 key 对应的数据状态不包含在返回的结果集中。

`GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)`：根据指定的 key 查询所有的历史记录信息。

`CreateCompositeKey(objectType string, attributes []string) (string, error)`：创建一个复合键。

`SplitCompositeKey(compositeKey string) (string, []string, error)`：将指定的复合键进行分割。

`GetQueryResult(query string) (StateQueryIteratorInterface, error)`：对(支持富查询功能的)状态数据库进行富查询，目前支持富查询的只有 CouchDB。

### 6.2.3 交易信息相关API
`GetTxID() string`：返回交易提案中指定的交易 ID。

`GetChannelID() string`：返回交易提案中指定的 Channel ID。

`GetTxTimestamp() (*timestamp.Timestamp, error)`：返回交易创建的时间戳，这个时间戳是peer 接收到交易的具体时间。

`GetBinding() ([]byte, error)`：返回交易的绑定信息。如果一些临时信息，以避免重复性攻击。

`GetSignedProposal() (*pb.SignedProposal, error)`：返回与交易提案相关的签名身份信息。

`GetCreator() ([]byte, error)`：返回该交易提交者的身份信息。

`GetTransient() (map[string][]byte, error)`：返回交易中不会被写至账本中的一些临时信息。

### 6.2.4 事件处理API
`SetEvent(name string, payload []byte) error`：设置事件，包括事件名称及内容。

### 6.2.5 对 PrivateData 操作的 API
`GetPrivateData(collection, key string) ([]byte, error)`：根据指定的 key，从指定的私有数据集中查询对应的私有数据。

`PutPrivateData(collection string, key string, value []byte) error`：将指定的 key 与 value 保存到私有数据集中。

`DelPrivateData(collection, key string) error`：根据指定的 key 从私有数据集中删除相应的数据。

`GetPrivateDataByRange(collection, startKey, endKey string) (StateQueryIteratorInterface, error)`：根据指定的开始与结束 key 查询范围（不包含结束key）内的私有数据。

`GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (StateQueryIteratorInterface, error)`：根据给定的部分组合键的集合，查询给定的私有状态。

`GetPrivateDataQueryResult(collection, query string) (StateQueryIteratorInterface, error)`：根据指定的查询字符串执行富查询 （只支持支持富查询的 CouchDB）。

## 6.3 链码实现的Hello World
> 前面我们已经接触了与链码相关的内容，下面我们根据已掌握的链码知识实现一个简单的链码应用。该应用需求较为简单：链码在实例化时向账本保存一个初始数据，key 为 Hello， value 为 World，然后用户发出查询请求，可以根据 key 查询到相应的 value。

### 6.3.1 链码开发
1. 创建文件夹  
进入 fabric-samples/chaincode/ 目录下并创建一个名为 hello 的文件夹
    ```shell
    cd hyfa/fabric-samples/chaincode
    sudo mkdir hello
    cd hello
    ```

2. 创建并编辑链码文件
    ```shell
    sudo vim hello.go
    ```

3. 导入链码依赖包
    ```go
    package main

    import (
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
    "fmt"
    )
    ```

4. 编写主函数
    ```shell
    func main()  {
    err := shim.Start(new(HelloChaincode))
    if err != nil {
        fmt.Printf("链码启动失败: %v", err)
    }
    }
    ```
5. 自定义结构体
    ```go
    type HelloChaincode struct {

    }
    ```
6. 实现 Chaincode 接口  
实现 Chaincode 接口必须重写 Init 与 Invoke 两个方法。  
**Init 函数：初始化数据状态**
   * 获取参数并判断参数长度是否为2
   * 参数: Key, Value
   * 调用 PutState 函数将状态写入账本中
   * 如果有错误, 则返回
   * 打印输出提示信息
   * 返回成功  

    具体实现代码如下：
    ```go
    // 实例化/升级链码时被自动调用
    // -c '{"Args":["Hello","World"]'
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
    ```
    **Invoke 函数**
   * 获取参数并判断长度是否为1
   * 利用第1个参数获取对应状态 GetState(key)
   * 如果有错误则返回
   * 如果返回值为空则返回错误
   * 返回成功状态  

    具体实现代码如下：
    ```golang
    // 对账本数据进行操作时被自动调用(query, invoke)
    func (t *HelloChaincode)  Invoke(stub shim.ChaincodeStubInterface) peer.Response  {
        // 获取调用链码时传递的参数内容(包括要调用的函数名及参数)
        fun, args := stub.GetFunctionAndParameters()

        // 客户意图
        if fun == "query"{
            return query(stub, args)
        }

            return shim.Error("非法操作, 指定功能不能实现")
    }
    ```
     实现查询函数

    函数名称为 query，具体实现如下：
    ```golang
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
    ```

### 6.3.2 链码测试
1. 启动网络  
进入 fabric-samples/chaincode-docker-devmode/ 目录
    ```shell
    cd ../chaincode-docker-devmode/
    ```
2. 构建并启动链码
    * 2.1 打开一个新的终端2，进入 chaincode 容器：
        ```shell
        sudo docker exec -it chaincode bash
        ```
    * 2.2 编译链码
        ```shell
        cd hello
        go build
        ```
    * 2.3 启动链码
        ```shell
        CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=hellocc:0 ./hello
        ```
        命令执行后终端输出如下：
        ```shell
        [shim] SetupChaincodeLogging -> INFO 001 Chaincode log level not provided; defaulting to: INFO
        [shim] SetupChaincodeLogging -> INFO 002 Chaincode (build level: ) starting up ...
        ```
3. 测试：
    * 3.1 打开一个新的终端3，进入 cli 容器
        ```shell
        sudo docker exec -it cli bash
        ```
    * 3.2 安装链码
        ```shell
        peer chaincode install -p chaincodedev/chaincode/hello -n hellocc -v 0
        ```
    * 3.3 实例化链码
        ```shell
        peer chaincode instantiate -n hellocc -v 0 -c '{"Args":["init", "Hello","World"]}' -C mycc
        ```  
    * 3.4 调用链码  
        根据指定的 key （"Hello"）查询对应的状态数据:
        ```shell
        peer chaincode query -n hellocc  -c '{"Args":["query","Hello"]}' -C mycc
        ```
        返回查询结果： World


## 6.4 动手编码一：链码实现资产管理
> 下面我们来实现一个简单的资产链码应用，该链码能够让用户在分类账上创建资产，并通过指定的函数实现对资产的修改与查询功能。

### 6.4.1 资产链码开发
1. 创建目录  
为 chaincode 应用创建一个名为 test 的目录
    ```shell
    cd ~/hyfa/fabric-samples/chaincode
    sudo mkdir test 
    cd test
    ```

2. 新建并编辑链码文件  
新建一个文件 test.go ，用于编写Go代码(这里用gedit编辑器，换成其他你顺手的编辑器都行)
    ```
    sudo gedit test.go
    ```
3. 导入链码依赖包  
    ```go
    package main

    import (
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
    "fmt"
    )
    ```
4. 定义结构体
    ```go
    type SimpleChaincode struct {
    }
    ```
5. 编写main函数
    ```go
    func main(){
        err := shim.Start(new(SimpleChaincode))
        if err != nil{
            fmt.Printf("启动 SimpleChaincode 时发生错误: %s", err)
        }
    }
    ```
6. 实现 Chaincode 接口  
Init 函数：初始化数据状态:  
* 获取参数, 使用 GetStringArgs 函数传递给调用链码的所需参数
* 检查合法性, 检查参数数量是否为2个, 如果不是, 则返回错误信息
* 利用两个参数, 调用 PutState 方法向账本中写入状态, 如果有错误则返回 shim.Error()， 否则返回 nil（shim.Success）
具体实现代码如下：
    ```go
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
    ```
    Invoke函数：验证函数名称为 set 或 get，并调用那些链式代码应用程序函数，通过 shim.Success 或 shim.Error 函数返回响应。
* 获取函数名与参数
* 对获取到的参数名称进行判断, 如果为 set, 则调用 set 方法, 反之调用 get
* set/get 函数返回两个值（result, err）
* 如果 err 不为空则返回错误
* err 为空则返回 []byte（result）  
具体实现代码如下：
```go
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
```
7. 实现具体业务功能的函数  
应用程序实现了两个可以通过 Invoke 函数调用的函数 （set/get）  
为了访问分类账的状态，利用 chaincode shim API 的 ChaincodeStubInterface.PutState 和ChaincodeStubInterface.GetState 函数  

    **7.1 实现set函数：修改资产**
   * 检查参数个数是否为2
   * 利用 PutState 方法将状态写入
   * 如果成功,则返回要写入的状态, 失败返回错误: fmt.Errorf("...")  
    具体实现代码如下：
    ```go
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
    ```

    **7.2 实现get函数：查询资产**  
    * 接收参数并判断个数 是否为1个
    * 调用 GetState 方法返回并接收两个返回值（value, err）判断 err 及 value 是否为空 return ""， fmt.Errorf("......")
    * 返回值 return string(value)，nil  
    具体实现代码如下：
        ```go
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
        ```
### 6.4.2 链码测试
使用 cd 跳转至 fabric-samples 的 chaincode-docker-devmode 目录

1. 终端1 启动网络
    ```shell
    sudo docker-compose -f docker-compose-simple.yaml up -d
    ```
    >在执行启动网络的命令之前确保无 Fabric 网络处于运行状态，如果有网络在运行，请先关闭。

2. 终端2 建立并启动链码
    * 2.1 打开一个新终端2，进入 chaincode 容器
        ```shell
        sudo docker exec -it chaincode bash
        ```
    * 2.2 编译   
    进入 test 目录编译 chaincode
        ```shell
        cd test
        go build
        ```
    * 2.3 运行chaincode
        ```shell
        CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=test:0 ./test
        ```
        命令执行后输出如下:
        ```shell
        [shim] SetupChaincodeLogging -> INFO 001 Chaincode log level not provided; defaulting to: INFO
        [shim] SetupChaincodeLogging -> INFO 002 Chaincode (build level: ) starting up ...
        ```
3. 终端3 测试  
    * 3.1 打开一个新的终端3，进入 cli 容器
        ```shell
        sudo docker exec -it cli bash
        ```
    * 3.2 安装链码
        ```shell
        peer chaincode install -p chaincodedev/chaincode/test -n test -v 0
        ```
    * 3.3 实例化链码
        ```shell
        peer chaincode instantiate -n test -v 0 -c '{"Args":["a","10"]}' -C mycc
        ```
    * 3.4 调用链码  
    指定调用 set 函数，将a的值更改为20
        ```shell
        peer chaincode invoke -n test -c '{"Args":["set", "a", "20"]}' -C mycc
        ```
        执行成功，输出如下内容：
        ```shell
        ......
        [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 0a8 Chaincode invoke successful. result: status:200 payload:"a"
        ```
    * 3.5 查询  
    指定调用 get 函数，查询 a 的值
        ```shell
        peer chaincode query -n test -c '{"Args":["query","a"]}' -C mycc
        ```
        执行成功, 输出: 20



## 6.5 动手编码二：链码实现转账
> 下面我们来实现一个使用链码能够实现对账户的查询，转账，删除账户的功能，并且整合完善资产管理应用链码的功能，该链码能够让用户在分类账上创建资产，并通过指定的函数实现对资产的修改与查询。

### 6.5.1 转账链码开发
1. 创建目录  
为 chaincode 应用创建一个名为 payment 的目录
    ```shell
    cd ~/hyfa/fabric-samples/chaincode
    sudo mkdir payment 
    cd payment
    ```
2. 新建并编辑链码文件  
新建一个文件 payment.go ，用于编写Go代码
    ```shell
    sudo vim payment.go
    ```
3. 导入链码依赖包
    ```go 
    package main

    import (
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
    "fmt"
    "strconv"
    )
    ```
4. 定义结构体
    ```go
    type PaymentChaincode struct {

    }
    ```
5. 编写main函数
    ```go
    func main(){
        err := shim.Start(new(PaymentChaincode))
        if err != nil{
            fmt.Printf("启动 PaymentChaincode 时发生错误: %s", err)
        }
    }
    ```
6. 实现 Chaincode 接口  
    **Init 函数：初始化两个账户，账户名分别为a、b，对应的金额为 100、200**  
   * 判断参数个数是否为4
   * 获取 args[0] 的值赋给A
   * strconv.Atoi（args[1]） 转换为整数, 返回 aval, err
   * 判断 err
   * 获取 args[2] 的值赋给B
   * strconv.Atoi（args[3]） 转换为整数, 返回 bval, err
   * 判断 err
   * 将 A 的状态值记录到分布式账本中
   * 判断 err
   * 将 B 的状态值记录到分布式账本中
   * 判断 err
   * return shim.Success（nil）  
    具体实现代码如下：
        ```go
        // 初始化两个账户及相应的余额
        // -c '{"Args":["init", "第一个账户名称", "第一个账户初始余额", "第二个账户名称", "第二个账户初始余额"]}'
        func (t *PaymentChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

            // 获取参数并验证
            _, args := stub.GetFunctionAndParameters()
            if len(args) != 4 {
                return shim.Error("必须指定两个账户名称及相应的初始余额")
            }

            // 判断账户名称是否合法
            var a = args[0]
            var avalStr = args[1]
            var b = args[2]
            var bvalStr = args[3]

            if len(a) < 2 {
                return shim.Error(a + " 账户名称不能少于2个字符长度")
            }
            if len(b) < 2 {
                return shim.Error(b + " 账户名称不能少于2个字符长度")
            }

            _, err := strconv.Atoi(avalStr)
            if err != nil {
                return shim.Error("指定的账户初始余额错误: " + avalStr)
            }
            _, err = strconv.Atoi(bvalStr)
            if err != nil {
                return shim.Error("指定的账户初始余额错误: " + bvalStr)
            }

            // 保存两个账户状态至账本中
            err = stub.PutState(a, []byte(avalStr))
            if err != nil {
                return shim.Error(a + " 保存状态时发生错误")
            }
            err = stub.PutState(b, []byte(bvalStr))
            if err != nil {
                return shim.Error(b + " 保存状态时发生错误")
            }

            return shim.Success([]byte("初始化成功"))

        }
        ```
    **Invoke 函数**：应用程序将具有三个不同的分支功能：find 、payment 、delete分别实现转账、删除、查询的功能, 根据交易参数定位到不同的分支处理逻辑。
   * 获取函数名称与参数列表
   * 判断函数名称并调用相应的函数  
    具体实现代码如下：
        ```go
        // peer chaincode query -n pay -C mycc -c '{"Args":["find", "a"]}'
        func (t *PaymentChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
        // 获取用户意图
        fun, args := stub.GetFunctionAndParameters()

        if fun == "find" {
            return find(stub, args)
        }else if fun == "payment" {
            return payment(stub, args)
        }else if fun == "del" {
            return delAccount(stub, args)
        }else if fun == "set" {
            return t.set(stub, args)
        }else if fun == "get" {
            return t.get(stub, args)
        }

        return shim.Error("非法操作, 指定的功能不能实现")
        }
        ```
7. 实现具体业务功能的函数  
应用程序实现了三个可以通过 Invoke 函数调用的函数（finde、payment、delAccount）  
    **7.1 实现 find 函数：根据给定的账户名称查询对应的状态信息**  
      * 判断参数是否为1个
      * 根据传入的参数调用 GetState 查询状态， aval， err 为接收返回值
      * 如果返回 err 不为空，则返回错误
      * 如果返回的状态为空，则返回错误
      * 如果无错误，返回查询到的值  
        具体实现代码如下：
        ```go
        // 根据指定的账户名称查询对应的余额信息
        // -c '{"Args":["find", "账户名称"]}'
        func find(stub shim.ChaincodeStubInterface, args []string) peer.Response {
            if len(args) != 1 {
                return shim.Error("必须且只能指定要查询的账户名称")
            }

            result, err := stub.GetState(args[0])
            if err != nil {
                return shim.Error("查询 " + args[0] + " 账户信息失败" + err.Error())
            }

            if result == nil {
                return shim.Error("根据指定 " + args[0] + " 没有查询到对应的余额")
            }

            return shim.Success(result)

        }
        ```
    **7.2 实现 payment 函数：根据指定的两个账户名称及金额，实现转账**
      * 判断参数是否为3
      * 获取两个账户名称（args[0] 与 args[1]）值, 赋给两个变量
      * 调用 GetState 获取 a 账户状态，avalsByte， err 为返回值
      * 判断有无错误（err 不为空， avalsByte 为空）
      * 类型转换: aval， _ = strconv.Atoi（string(avalsByte)）
      * 调用 GetState 获取 b 账户状态， bvalsByte，err 为返回值
      * 判断有无错误（err 不为空，bvalsByte 为空）
      * 类型转换: bval， _ = strconv.Atoi（string(bvalsByte)）
      * 将要转账的数额进行类型转换： x， err = strconv.Atoi（args[2]）
      * 判断 err 是否为空
      * aval， bval 执行转账操作
      * 记录状态， err = PutState(a, []byte（strconv.Itoa(aval))）
      * Itoa： 将整数转换为十进制字符串形式
      * 判断有无错误.
      * 记录状态， err = PutState（b, []byte(strconv.Itoa(bval))）
      * 判断有无错误.
      * return shim.Success（nil）  
      具体实现代码如下：
        ```go
        // 转账
        // -c '{"Args":["payment", "源账户名称", "目标账户名称", "转账金额"]}'
        func payment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
            if len(args) != 3 {
                return shim.Error("必须且只能指定源账户及目标账户名称与对应的转账金额")
            }

            var source, target string
            var x string

            source = args[0]
            target = args[1]
            x = args[2]

            // 源账户扣除对应的转账金额
            // 目标账户加上对应的转账金额

            // 查询源账户及目标账户的余额
            sval, err := stub.GetState(source)
            if err != nil {
                return shim.Error("查询源账户信息失败")
            }
            // 如果源账户或目标账户不存在的情况下
            // 不存在的情况下直接return

            tval, err := stub.GetState(target)
            if err != nil {
                return shim.Error("查询目标账户信息失败")
            }

            // 实现转账
            s, err := strconv.Atoi(x)
            if err != nil {
                return shim.Error("指定的转账金额错误")
            }

            svi, err := strconv.Atoi(string(sval))
            if err != nil {
                return shim.Error("处理源账户余额时发生错误")
            }

            tvi, err := strconv.Atoi(string(tval))
            if err != nil {
                return shim.Error("处理目标账户余额时发生错误")
            }

            if svi < s {
                return shim.Error("指定的源账户余额不足, 无法实现转账")
            }

            svi = svi - s
            tvi = tvi + s

            // 将修改之后的源账户与目标账户的状态保存至账本中
            err = stub.PutState(source, []byte(strconv.Itoa(svi)))
            if err != nil {
                return  shim.Error("保存转账后的源账户状态失败")
            }

            err = stub.PutState(target, []byte(strconv.Itoa(tvi)))
            if err != nil {
                return  shim.Error("保存转账后的目标账户状态失败")
            }

            return shim.Success([]byte("转账成功"))

        }
        ```
   **7.3 实现 delAccount 函数：根据指定的名称删除对应的实体信息**  
      * 判断参数个数是否为1
      * 调用 DelState 方法，err 接收返回值
      * 如果 err 不为空, 返回错误
      * 返回成功 shim.Success(nil)  
        具体实现代码如下：
        ```go
        // 根据指定的账户名称删除相应信息
        // -c '{"Args":["del", "账户名称"]}'
        func delAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
        if len(args) != 1 {
            return shim.Error("必须且只能指定要删除的账户名称")
        }

        result, err := stub.GetState(args[0])
        if err != nil {
            return shim.Error("查询 " + args[0] + " 账户信息失败" + err.Error())
        }

        if result == nil {
            return shim.Error("根据指定 " + args[0] + " 没有查询到对应的余额")
        }

        err = stub.DelState(args[0])
        if err != nil {
            return shim.Error("删除指定的账户失败: " + args[0] + ", " + err.Error())
        }

        return shim.Success([]byte("删除指定的账户成功" + args[0]))
        }
        ```
    **7.4 实现 set 函数，设置指定账户的值**  
在简单资产管理链码的的 set 函数的功能并不完善，因为我们没有考虑用户存入资产之后需要对该账户的资产进行修改，现在我们来添加这一功能。  
    具体实现代码如下：
    ```go 
    // 向指定的账户存入对应的金额
    // -c '{"Args":["set", "账户名称", "要存入的金额"]}'
    func (t *PaymentChaincode) set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
        if len(args) != 2 {
            return shim.Error("必须且只能指定账户名称及要存入的金额")
        }

        result, err := stub.GetState(args[0])
        if err != nil {
            return shim.Error("根据指定的账户查询信息失败")
        }

        if result == nil {
            return shim.Error("指定的账户不存在")
        }

        // 存入账户
        val, err := strconv.Atoi(string(result))
        if err != nil {
            return shim.Error("处理指定的账户金额时发生错误")
        }
        x, err := strconv.Atoi(args[1])
        if err != nil {
            return shim.Error("指定要存入的金额错误")
        }

        val = val + x

        // 保存信息
        err = stub.PutState(args[0], []byte(strconv.Itoa(val)))
        if err != nil {
            return shim.Error("存入账户金额时发生错误")
        }
        return shim.Success([]byte("存入操作成功"))

    }
    ```
    **7.5 实现 get 函数，从指定的账户中提取指定的金额**  
    同理，用户从账户中提取从指定金额的资产之后，也需要对该账户的资产进行修改。   
    具体实现代码如下：
    ```go
    // 从账户中提取指定的金额
    // -c '{"Args":["get", "账户名称", "要提取的金额"]}'
    func (t *PaymentChaincode) get(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
        if len(args) != 2 {
            return shim.Error("必须且只能指定要提取的账户名称及金额")
        }

        x, err := strconv.Atoi(args[1])
        if err != nil {
            return shim.Error("指定要提取的金额错误, 请重新输入")
        }

        // 从指定的账户中查询出现有金额
        result, err := stub.GetState(args[0])
        if err != nil {
            return shim.Error("查询指定账户金额时发生错误")
        }
        if result == nil {
            return shim.Error("要查询的账户不存在或已被注销")
        }

        val, err := strconv.Atoi(string(result))
        if err != nil {
            return shim.Error("处理账户金额时发生错误")
        }

        if val < x {
            return shim.Error("要提取的金额不足")
        }

        val = val - x
        err = stub.PutState(args[0], []byte(strconv.Itoa(val)))
        if err != nil {
            return shim.Error("提取失败, 保存数据时发生错误")
        }
        return shim.Success([]byte("提取成功"))

    }
    ```
### 6.5.2 链码测试
跳转至 fabric-samples 的 chaincode-docker-devmode 目录

1. 终端1 启动网络
    ```shell
    sudo docker-compose -f docker-compose-simple.yaml up -d
    ```
    > 在执行启动网络的命令之前确保无Fabric网络处于运行状态，如果有网络在运行，请先关闭。

2. 终端2 建立并启动链码
    * 2.1 打开一个新终端2，进入 chaincode 容器
        ```shell
        sudo docker exec -it chaincode bash
        ```
    * 2.2 编译  
        进入 test 目录编译 chaincode
        ```shell
        cd payment
        go build
        ```
    * 2.3 运行chaincode
        ```shell
        CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=paycc:0 ./payment
        ```
        命令执行后输出如下:
        ```shell
        [shim] SetupChaincodeLogging -> INFO 001 Chaincode log level not provided; defaulting to: INFO
        [shim] SetupChaincodeLogging -> INFO 002 Chaincode (build level: ) starting up ...
        ```
3. 终端3 测试
    * 3.1 打开一个新的终端3，进入 cli 容器
        ```shell
        sudo docker exec -it cli bash
        ```
    * 3.2 安装链码
        ```shell
        peer chaincode install -p chaincodedev/chaincode/payment -n paycc -v 0
        ```
    * 3.3 实例化链码
        ```shell
        peer chaincode instantiate -n paycc -v 0 -c '{"Args":["init","aaa", "100", "bbb","200"]}' -C mycc
        ```
    * 3.4 调用链码  
    指定调用 payment 函数，从 aaa 账户向 bbb 账户转账 20
        ```shell
        peer chaincode invoke -n paycc -c '{"Args":["payment", "aaa","bbb","20"]}' -C mycc
        ```
        执行成功，输出如下内容：
        ```shell
        ......
        [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 0a8 Chaincode invoke successful. result: status:200 payload:"\350\275\254\350\264\246\346\210\220\345\212\237"
        ```
    * 3.5 查询  
    * 指定调用 find 函数，查询 a 账户的值
        ```shell
        peer chaincode query -n paycc -c '{"Args":["find","aaa"]}' -C mycc
        ```
        执行成功, 输出: 80