# Course 05

## contents

### 5.认识链码

我们已经学习了通过一系列命令开启一个fabric的网络，了解了配置文件如何设定网络的架构信息，实践了通过修改配置文件设计自己的网络，接下来我将教大家
链码的知识和fabric上链码的管理。

#### 5.1 链码的概念

Chaincode是一段由Go语言编写的（支持其他编程语言，如Java），并能实现预定义接口的程序。Chaincode运行在一个受保护的Docker容器当中，与背书堆栈的运行相互隔离。Chaincode可通过应用提交的交易对账本状态初始化并进行管理。

一段链码通常处理由网络中的成员一致认可的业务逻辑，故我们很可能用“智能合约”来代指链码。一段chiancode创建的（账本）状态是与其他链码相互隔离的，故而不能被其他链码直接访问。不过，如果是在相同的网络中，一段chiancode在获取相应许可后则可以调用其他chiancode来访问它的账本。

<table border=0 cellpadding=0 cellspacing=0 width=684 style='border-collapse:
 collapse;table-layout:fixed;width:514pt'>
 <col width=241 style='mso-width-source:userset;mso-width-alt:8424;width:181pt'>
 <col width=251 style='mso-width-source:userset;mso-width-alt:8773;width:189pt'>
 <col width=192 style='mso-width-source:userset;mso-width-alt:6702;width:144pt'>
 <tr height=31 style='height:23.0pt'>
  <td colspan=2 height=31 class=xl66 width=492 style='height:23.0pt;width:370pt'>系统合约</td>
  <td class=xl73 width=192 style='width:104pt'>用户合约</td>
 </tr>
 <tr height=29 style='height:21.5pt'>
  <td height=29 class=xl67 dir=LTR width=241 style='height:21.5pt;width:181pt'>配置系统链码<font
  class="font8">(CSCC)</font></td>
  <td class=xl68 dir=LTR width=251 style='border-left:none;width:189pt'>Peer <font
  class="font10">端的 </font><font class="font9">Channel </font><font
  class="font10">配置</font></td>
  <td rowspan=5 class=xl74 width=192 style='width:144pt'>"由应用程序开发人员根据不同场景需求及成员制定的相关规则，使用 Golang或Java等语言编写的基于操作区块链分布式账本的状态的业务处理逻辑代码，运行在链码容器中，通过Fabric 提供的接口与账本状态进行交互。下可对账本数据进行操作，上可以给企业级应用程序提供调用接口<br>
    </td>
 </tr>
 <tr height=52 style='height:39.0pt'>
  <td height=52 class=xl69 dir=LTR width=241 style='height:39.0pt;border-top:
  none;width:181pt'>生命周期系统链码<font class="font12">(LSCC)</font></td>
  <td class=xl70 dir=LTR width=251 style='border-top:none;border-left:none;
  width:189pt'>对用户链码的生命周期进行管理</td>
 </tr>
 <tr height=54 style='height:40.5pt'>
  <td height=54 class=xl71 dir=LTR width=241 style='height:40.5pt;border-top:
  none;width:181pt'>查询系统链码<font class="font12">(QSCC)</font></td>
  <td class=xl72 dir=LTR width=251 style='border-top:none;border-left:none;
  width:189pt'>提供账本查询<font class="font14">API</font><font class="font13">，如获取区块和交易等信息</font></td>
 </tr>
 <tr height=54 style='height:40.5pt'>
  <td height=54 class=xl71 dir=LTR width=241 style='height:40.5pt;border-top:
  none;width:181pt'>背书管理系统链码<font class="font12">(ESCC)</font></td>
  <td class=xl72 dir=LTR width=251 style='border-top:none;border-left:none;
  width:189pt'>负责背书(签名)过程, 并可以支持对背书策略进行管理</td>
 </tr>
 <tr height=51 style='height:38.5pt'>
  <td height=51 class=xl71 dir=LTR width=241 style='height:38.5pt;border-top:
  none;width:181pt'>验证系统链码<font class="font12">(VSCC)</font></td>
  <td class=xl72 dir=LTR width=251 style='border-top:none;border-left:none;
  width:189pt'>处理交易的验证，包括检查背书策略以及多版本并发控制</td>
 </tr>
 <![if supportMisalignedColumns]>
 <tr height=0 style='display:none'>
  <td width=241 style='width:181pt'></td>
  <td width=251 style='width:189pt'></td>
  <td width=192 style='width:144pt'></td>
 </tr>
 <![endif]>
</table>

管理 Chaincode 的五个命令：

+ install：将已编写完成的链码安装在网络节点中。

+ instantiate：对已安装的链码进行实例化。

+ upgrade：对已有链码进行升级。链代码可以在安装后根据具体需求的变化进行升级。

+ package：对指定的链码进行打包的操作。

+ singnpackage：签名。
#### 5.2 链码的欣赏
[github.com/chaincode/chaincode_example02/go/](files/chaincode_example02.go)

[chaincode基础架构](files/chaincode.go)
#### 5.3 链码的管理
我们对链码已经有了一个基础的认识，下面我们利用 fabric-samples 提供的示例链码来进行实践；如何安装、实例化、调用、打包、签名和升级链码等操作。

首先确认网络是否处于开启状态，利用 docker ps 命令查看容器是否处于活动状态
```
$ sudo docker ps
```
如果没有活动的容器，则先使用 docker-compose 命令启动网络然后进入CLI 容器中
```
$ sudo docker-compose -f docker-compose-cli.yaml up -d
$ sudo docker exec -it cli bash
```
如果当前已进入至 CLI 容器中，则上面的命令无需执行。如果之前使用 exit 命令退出了 cli 容器，请使用 
```
sudo docker exec -it cli bash 
```
命令重新进入 cli 容器。

检查当前节点（默认为peer0.example.com）已加入到哪些通道中：
```
# peer channel list
```
执行成功后会在终端中输出：
```
Channels peers has joined: 
```
mychannel
根据如下的输出内容，说明当前节点已成功加入到一个名为 mychannel 的应用通道中。Peer加入应用通道后，可以执行链码调用的相关操作，进行测试。如果没有，则先将当前节点加入到已创建的应用通道中。

检查环境变量是否正确设置
```
# echo $CHANNEL_NAME
```
设置环境变量，指定应用通道名称为 mychannel ，因为我们创建的应用通道及当前的 peer 节点加入的应用通道名称为 mychannel
```
# export CHANNEL_NAME=mychannel
```
链码调用处理交易之前必须将其部署到 Peer 节点上，实现步骤如下：

将其安装在指定的网络节点上
安装完成后要对其进行实例化
然后才可以调用链码处理交易(查询或执行事务)
##### 5.3.1 链码的安装
安装链码使用 install 命令：
```
# peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/chaincode_example02/go/
```
参数说明：
```
-n： 指定要安装的链码的名称
-v： 指定链码的版本
-p： 指定要安装的链码的所在路径
```
命令执行完成看到如下输出说明指定的链码被成功安装至 peer 节点：
```
[chaincodeCmd] checkChaincodeCmdParams -> INFO 001 Using default escc
[chaincodeCmd] checkChaincodeCmdParams -> INFO 002 Using default vscc
[chaincodeCmd] install -> INFO 003 Installed remotely response:<status:200 payload:"OK" >
```


##### 5.3.2 链码的实例化
实例化链码使用 instantiate 命令：
```
# peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
```
参数说明:
```
-o： 指定Oderer服务节点地址
--tls： 开启 TLS 验证
--cafile： 指定 TLS_CA 证书的所在路径
-n： 指定要实例化的链码名称，必须与安装时指定的链码名称相同
-v： 指定要实例化的链码的版本号，必须与安装时指定的链码版本号相同
-C： 指定通道名称
-c： 实例化链码时指定的参数
-P： 指定背书策略
```
实例化完成后，用户即可向网络中发起交易。


##### 5.3.3 链码的调用
调用查询函数
使用 query 命令实现：
```
# peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
```
参数说明：
```
-n： 指定要调用的链码名称
-C： 指定通道名称
-c 指定调用链码时所需要的参数
执行成功输出结果：100
```

##### 调用invoke函数
调用链码使用 invoke 命令实现：
```
# peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C $CHANNEL_NAME -n mycc -c '{"Args":["invoke","a","b","10"]}'
```
参数说明：
```
-o： 指定orderer节点地址
--tls： 开启TLS验证
--cafile： 指定TLS_CA证书路径
-n: 指定链码名称
-C： 指定通道名称
-c： 指定调用链码的所需参数
```
有如下输出则说明链码被调用成功且交易请求被成功处理：
```
[chaincodeCmd] chaincodeInvokeOrQuery -> INFO 001 Chaincode invoke successful. result: status:200
```
##### 查询a账户的金额
执行查询a账户的命令，并查看输出结果：
```
# peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
```
执行成功输出结果: 90

##### 5.3.4 链码的打包
通过将链码相关数据（如链码名称、版本、实例化策略等信息）进行封装，可以实现对其进行打包和签名的操作。

chaincode 包具体包含以下三个部分：

+ chaincode 本身，由 ChaincodeDeploymentSpec（CDS）定义。CDS 根据代码及一些其他属性（名称，版本等）来定义 chaincode。
+ 一个可选的实例化策略，该策略可被 背书策略 描述。
+ 一组表示 chaincode 所有权的签名。

对于一个已经编写完成的链码可以使用 package 命令进行打包操作：
```
# peer chaincode package -n exacc -v 1.0 -p github.com/chaincode/chaincode_example02/go/  -s -S -i "AND('Org1MSP.admin')" ccpack.out
```
参数说明：
```
-s： 创建一个可以被多个所有者签名的包。

-S： 可选参数，使用 core.yaml 文件中被 localMspId 相关属性值定义的 MSP 对包进行签名。

-i： 指定链码的实例化策略（指定谁可以实例化链码）。
```
打包后的文件，可以直接用于 install 操作，如： peer chaincode install ccpack.out，但我们一般会在将打包的文件进行签名之后再做进一步的处理。
##### 5.3.5 链码的签名
对一个打包文件进行签名操作（添加当前 MSP 签名到签名列表中）

使用 signpackage 命令实现：
```
# peer chaincode signpackage ccpack.out signedccpack.out
```
signedccpack.out 包含一个用本地 MSP 对包进行的附加签名。

添加了签名的链码包可以进行进行一步的处理，如先将链码进行安装，然后对已安装的链码进行实例化或升级的操作。

安装已添加签名的链码：
```
# peer chaincode install signedccpack.out
```
命令执行成功输出如下内容：
```
[chaincodeCmd] install -> INFO 001 Installed remotely response:<status:200 payload:"OK" >
```
安装成功之后进行链码的实例化操作，同时指定其背书策略。
```
# peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n exacc -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
```

##### 5.3.6 链码的升级
在实际场景中，由于需求场景的变化，链码也需求实时做出修改，以适应不同的场景需求。所以我们必须能够对于已成功部署并运行状态中的链码进行升级操作。

首先，先将修改之后的链码进行安装，然后使用 upgrade 命令对已安装的链码进行升级，具体实现如下：

安装：
```
# peer chaincode install -n mycc -v 2.0 -p github.com/chaincode/chaincode_example02/go/
```
升级：
```
# peer chaincode upgrade -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc -v 2.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
```

测试方式和链码的调用是一样的：

查询链码：
```
# peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
```
执行成功输出查询结果： 100

调用链码：
```
# peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C $CHANNEL_NAME -n mycc -c '{"Args":["invoke","a","b","10"]}'
```
执行成功输出如下：
```
[chaincodeCmd] chaincodeInvokeOrQuery -> INFO 001 Chaincode invoke successful. result: status:200
```
查询链码：
```
# peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
```
执行成功输出查询结果： 90

#### 5.2 开发者模式下的链码开发
在 dev 开发模式下我们可以使用三个终端来实现具体的测试过程

##### 启动网络
终端1（当前终端）
关闭之前已启动的网络环境：
```
$ sudo docker-compose -f docker-compose-cli.yaml down
```
进入 chaincode-docker-devmode 目录
```
$ cd ~/hyfa/fabric-samples/chaincode-docker-devmode/
```

下面，我们使用 docker-compose-simple.yaml 配置文件来启动网络：
```
$ sudo docker-compose -f docker-compose-simple.yaml up -d
```
上面的命令以 docker-compose-simple.yaml 启动了网络，并以开发模式启动 peer。另外还启动了两个容器：

一个 chaincode 容器，用于链码环境

一个 CLI 容器，用于与链码进行交互。

命令执行后，终端中输出如下：
```
Creating orderer
Creating peer
Creating chaincode
Creating cli
```
创建和连接通道的命令嵌入到 CLI 容器中，因此我们可以立即跳转到链码调用。

##### 构建并启动链码
网络启动成功后，下一步需要开发者自行对已经编写好的链码进行构建及启动。

终端2（开启一个新的终端2）

###### 进入chaincode容器
chaincode 容器的作用是为了以简化的方式建立并启动链码
```
$ sudo docker exec -it chaincode bash
```
命令提示符变为：
```
root@858726aed16e:/opt/gopath/src/chaincode#
```
进入 chaincode 容器之后就可以构建与启动链码。

###### 编译链码
现在我们对 fabric-samples 提供的 chaincode_example02 进行测试，当然，在实际环境中，我们可以将开发的链码添加到 chaincode 子目录中并重新构建及启动链码，然后进行测试。

进入 chaincode_example02/go/ 目录编译 chaincode
```
# cd chaincode_example02/go/

# go build
```
###### 运行chaincode
使用如下命令启动并运行链码：
```
# CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./go
```
命令执行后输出如下：
```
[shim] SetupChaincodeLogging -> INFO 001 Chaincode log level not provided; defaulting to: INFO
[shim] SetupChaincodeLogging -> INFO 002 Chaincode (build level: ) starting up ...
```
命令含义：
```
CORE_PEER_ADDRESS：用于指定peer。
CORE_CHAINCODE_ID_NAME：用于注册到peer的链码。
mycc： 指定链码名称
0： 指定链码初始版本号
./go： 指定链码文件
```
注意，此阶段，链码与任何通道都没有关联。我们需要在后续步骤中使用“实例化”命令来完成.

##### 调用链码
终端3（开启一个新的终端3）

首先进入 cli 容器
```
$ sudo docker exec -it cli bash
```
进入 CLI 容器后，执行如下命令安装及实例化 chaincode

即使我们在 dev 模式下，也需要安装链码，使链码能够正常通过生命周期系统链码的检查 。将来可能会删除此步骤。

安装：
```
# peer chaincode install -p chaincodedev/chaincode/chaincode_example02/go -n mycc -v 0
```
注意：安装链码时指定的链码名称与版本号必须与在终端2中注册的链码名称及版本号相同。

安装命令执行后，终端中输出如下：
```
......
-----END CERTIFICATE-----
[msp] setupSigningIdentity -> DEBU 034 Signing identity expires at 2027-11-10 13:41:11 +0000 UTC
[msp] Validate -> DEBU 035 MSP DEFAULT validating identity
[grpc] Printf -> DEBU 036 parsed scheme: ""
[grpc] Printf -> DEBU 037 scheme "" not registered, fallback to default scheme
[grpc] Printf -> DEBU 038 ccResolverWrapper: sending new addresses to cc: [{peer:7051 0  <nil>}]
[grpc] Printf -> DEBU 039 ClientConn switching balancer to "pick_first"
[grpc] Printf -> DEBU 03a pickfirstBalancer: HandleSubConnStateChange: 0xc4204e7c40, CONNECTING
[grpc] Printf -> DEBU 03b pickfirstBalancer: HandleSubConnStateChange: 0xc4204e7c40, READY
[grpc] Printf -> DEBU 03c parsed scheme: ""
[grpc] Printf -> DEBU 03d scheme "" not registered, fallback to default scheme
[grpc] Printf -> DEBU 03e ccResolverWrapper: sending new addresses to cc: [{peer:7051 0  <nil>}]
[grpc] Printf -> DEBU 03f ClientConn switching balancer to "pick_first"
[grpc] Printf -> DEBU 040 pickfirstBalancer: HandleSubConnStateChange: 0xc420072170, CONNECTING
[grpc] Printf -> DEBU 041 pickfirstBalancer: HandleSubConnStateChange: 0xc420072170, READY
[msp] GetDefaultSigningIdentity -> DEBU 042 Obtaining default signing identity
[chaincodeCmd] checkChaincodeCmdParams -> INFO 043 Using default escc
[chaincodeCmd] checkChaincodeCmdParams -> INFO 044 Using default vscc
[chaincodeCmd] getChaincodeSpec -> DEBU 045 java chaincode disabled
[golang-platform] getCodeFromFS -> DEBU 046 getCodeFromFS chaincodedev/chaincode/chaincode_example02/go
[golang-platform] func1 -> DEBU 047 Discarding GOROOT package fmt
[golang-platform] func1 -> DEBU 048 Discarding provided package github.com/hyperledger/fabric/core/chaincode/shim
[golang-platform] func1 -> DEBU 049 Discarding provided package github.com/hyperledger/fabric/protos/peer
[golang-platform] func1 -> DEBU 04a Discarding GOROOT package strconv
[golang-platform] GetDeploymentPayload -> DEBU 04b done
[container] WriteFileToPackage -> DEBU 04c Writing file to tarball: src/chaincodedev/chaincode/chaincode_example02/go/chaincode_example02.go
[msp/identity] Sign -> DEBU 04d Sign: plaintext: 0AC4070A5C08031A0C08C3F492DC0510...21E3DF010000FFFF4C61C899001C0000 
[msp/identity] Sign -> DEBU 04e Sign: digest: 6F0F7CF70A07027506571AAC56B978353CA3C73E311C882AB57263543ECE7B76 
[chaincodeCmd] install -> INFO 04f Installed remotely response:<status:200 payload:"OK" >
```
实例化：
```
# peer chaincode instantiate -n mycc -v 0 -c '{"Args":["init","a", "100", "b","200"]}' -C myc
```
实例化命令执行后，终端中输出如下内容：
```
......
[common/configtx] addToMap -> DEBU 091 Adding to config map: [Policy] /Channel/Application/Readers
[common/configtx] addToMap -> DEBU 092 Adding to config map: [Policy] /Channel/Application/Writers
[common/configtx] addToMap -> DEBU 093 Adding to config map: [Policy] /Channel/Application/Admins
[common/configtx] addToMap -> DEBU 094 Adding to config map: [Value]  /Channel/BlockDataHashingStructure
[common/configtx] addToMap -> DEBU 095 Adding to config map: [Value]  /Channel/OrdererAddresses
[common/configtx] addToMap -> DEBU 096 Adding to config map: [Value]  /Channel/HashingAlgorithm
[common/configtx] addToMap -> DEBU 097 Adding to config map: [Value]  /Channel/Consortium
[common/configtx] addToMap -> DEBU 098 Adding to config map: [Policy] /Channel/Writers
[common/configtx] addToMap -> DEBU 099 Adding to config map: [Policy] /Channel/Admins
[common/configtx] addToMap -> DEBU 09a Adding to config map: [Policy] /Channel/Readers
[chaincodeCmd] InitCmdFactory -> INFO 09b Retrieved channel (myc) orderer endpoint: orderer:7050
[grpc] Printf -> DEBU 09c parsed scheme: ""
[grpc] Printf -> DEBU 09d scheme "" not registered, fallback to default scheme
[grpc] Printf -> DEBU 09e ccResolverWrapper: sending new addresses to cc: [{orderer:7050 0  <nil>}]
[grpc] Printf -> DEBU 09f ClientConn switching balancer to "pick_first"
[grpc] Printf -> DEBU 0a0 pickfirstBalancer: HandleSubConnStateChange: 0xc42043d790, CONNECTING
[grpc] Printf -> DEBU 0a1 pickfirstBalancer: HandleSubConnStateChange: 0xc42043d790, READY
[chaincodeCmd] checkChaincodeCmdParams -> INFO 0a2 Using default escc
[chaincodeCmd] checkChaincodeCmdParams -> INFO 0a3 Using default vscc
[chaincodeCmd] getChaincodeSpec -> DEBU 0a4 java chaincode disabled
[msp/identity] Sign -> DEBU 0a5 Sign: plaintext: 0AC9070A6108031A0C08F2F592DC0510...30300A000A04657363630A0476736363 
[msp/identity] Sign -> DEBU 0a6 Sign: digest: B7822DC27649C2CE85206E13DC69861CDB6C4786D6D3E299032BE2A187C0A362 
[msp/identity] Sign -> DEBU 0a7 Sign: plaintext: 0AC9070A6108031A0C08F2F592DC0510...025C39086D09D5D731F33C16A2E53492 
[msp/identity] Sign -> DEBU 0a8 Sign: digest: 27E503A393AD2B63F56A02FD29E4495999D913F037FEE4BCD894C16447EDAB35
```
测试还是和链码调用是一样的：
```
查询：

# peer chaincode query -n mycc  -c '{"Args":["query","a"]}' -C myc
执行成功输出查询结果： 100

调用：

# peer chaincode invoke -n mycc -c '{"Args":["invoke","a","b","10"]}' -C myc
执行成功输出如下：

[chaincodeCmd] chaincodeInvokeOrQuery -> INFO 0a8 Chaincode invoke successful. result: status:200
查询：

# peer chaincode query -n mycc  -c '{"Args":["query","a"]}' -C myc
执行成功输出查询结果： 90
```
