进入 chaincode 目录中，创建并进入 testcdb 目录：
```shell
cd hyfa/fabric-samples/chaincode
sudo mkdir testcdb
cd testcdb
```

将编写的 main.go 文件上传至 testcdb 目录中，然后跳转至fabric-samples的chaincode-docker-devmode目录

```shell
cd ~/hyfa/fabric-samples/chaincode-docker-devmode/
```

#### 1. 终端1 启动网络
```shell
sudo docker-compose -f docker-compose-simple.yaml up -d
```
> 在执行启动网络的命令之前确保无Fabric网络处于运行状态，如果有网络在运行，请先关闭。

#### 2. 终端2 建立并启动链码

**2.1 打开一个新终端2，进入 chaincode 容器**
```shell
sudo docker exec -it chaincode bash
```
**2.2 编译**  
进入 testcdb 目录编译 chaincode
```shell
cd testcdb
go build
```
**2.3 运行chaincode**
```shell
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=cdb:0 ./testcdb
```
命令执行后输出如下：
```shell
[shim] SetupChaincodeLogging -> INFO 001 Chaincode log level not provided; defaulting to: INFO
[shim] SetupChaincodeLogging -> INFO 002 Chaincode (build level: ) starting up ...
```
#### 3. 终端3 测试

**3.1 打开一个新的终端3，进入 cli 容器**
```shell
sudo docker exec -it cli bash
```
**3.2 安装链码**
```shell
peer chaincode install -p chaincodedev/chaincode/testcdb -n cdb -v 0
```
**3.3 实例化链码**
```shell
peer chaincode instantiate -n cdb -v 0 -C myc -c '{"Args":["init"]}'
```
**3.4 初始化数据**

指定调用 billInit 函数进行数据的初始化：
```shell
peer chaincode invoke -n cdb -C myc -c '{"Args":["carInit"]}'
```
执行成功，输出如下内容：
```shell
......
[chaincodeCmd] chaincodeInvokeOrQuery -> INFO 0a8 Chaincode invoke successful. result: status:200 payload:"\345\210\235\345\247\213\345\214\226\347\245\250\346\215\256\346\210\220\345\212\237"
```

**3.5 根据 owner 查询数据**

指定调用 queryCars 函数，查询指定用户 p1
```shell
peer chaincode query -n cdb -C myc -c '{"Args":["queryCars", "{\"selector\":{\"owner\":\"p1\"}}"]}'
```

执行成功，输出查询到的结果如下：
```shell
......
key:car001, Value: {"brand":"brand-1","carid":"car001","carname":"car-wife","docType":"carObj","owner":"p1","price":"20"}; key:car002, Value: {"brand":"brand-2","carid":"car002","carname":"car-me","docType":"carObj","owner":"p1","price":"40"}; key:car003, Value: {"brand":"audi-3","carid":"car003","carname":"car-son","docType":"carObj","owner":"p1","price":"20"}
```

**3.6 根据 owner 以及指定 brand 来查询**

```shell
peer chaincode query -n cdb -C myc -c '{"Args":["queryCars", "{\"selector\":{\"owner\":\"p1\", \"brand\":\"brand-2\"}}"]}'
```
结果：
```shell
key:car001, Value: {"brand":"brand-1","carid":"car001","carname":"car-wife","docType":"carObj","owner":"p1","price":"20"}; key:car002, Value: {"brand":"brand-2","carid":"car002","carname":"car-me","docType":"carObj","owner":"p1","price":"40"}
```

**3.7 根据 owner 以及 brand 的正则表达式匹配规则来查询**

指定调用 queryCars 函数：
```shell
peer chaincode query -n cdb -C myc -c '{"Args":["queryCars", "{\"selector\":{\"owner\":\"p1\", \"brand\":{\"$regex\":\"^brand\"}}}"]}'
```
执行成功，输出查询到的结果如下：
```shell
......
key:car001, Value: {"brand":"brand-1","carid":"car001","carname":"car-wife","docType":"carObj","owner":"p1","price":"20"}; key:car002, Value: {"brand":"brand-2","carid":"car002","carname":"car-me","docType":"carObj","owner":"p1","price":"40"}

```


最后注意清理网络。

参考：  
https://docs.couchdb.org/en/2.2.0/api/database/find.html

https://github.com/hyperledger/fabric-samples/blob/master/chaincode/marbles02/go/marbles_chaincode.go#L22