# Build Your First Network

## `Contents`
- [Build Your First Network](#build-your-first-network)
  - [`Contents`](#contents)
  - [1. 1.生成公私钥和证书](#1-1%e7%94%9f%e6%88%90%e5%85%ac%e7%a7%81%e9%92%a5%e5%92%8c%e8%af%81%e4%b9%a6)




## 1. 1.生成公私钥和证书
Fabric中有两种类型的公私钥和证书，一种是给节点之前通讯安全而准备的TLS证书，另一种是用户登录和权限控制的用户证书。这些证书本来应该是由CA来颁发，但是我们这里是测试环境，并没有启用CA节点，所以Fabric帮我们提供了一个工具：cryptogen。

1.1编译生成cryptogen
我们既然获得了Fabric的源代码，那么就可以轻易的使用make命令编译需要的程序。Fabric官方提供了专门编译cryptogen的入口，我们只需要运行以下命令即可：

cd ~/go/src/github.com/hyperledger/fabric
make cryptogen
运行后系统返回结果：

build/bin/cryptogen 
CGO_CFLAGS=" " GOBIN=/home/studyzy/go/src/github.com/hyperledger/fabric/build/bin go install -tags "" -ldflags "-X github.com/hyperledger/fabric/common/tools/cryptogen/metadata.Version=1.0.0" github.com/hyperledger/fabric/common/tools/cryptogen 
Binary available as build/bin/cryptogen
也就是说我们在build/bin文件夹下可以看到编译出来的cryptogen程序。

1.2配置crypto-config.yaml
examples/e2e_cli/crypto-config.yaml已经提供了一个Orderer Org和两个Peer Org的配置，该模板中也对字段进行了注释。我们可以把Org2拿来分析一下：

- Name: Org2 
  Domain: org2.example.com 
  Template: 
    Count: 2 
  Users: 
    Count: 1
Name和Domain就是关于这个组织的名字和域名，这主要是用于生成证书的时候，证书内会包含该信息。而Template Count=2是说我们要生成2套公私钥和证书，一套是peer0.org2的，还有一套是peer1.org2的。最后Users. Count=1是说每个Template下面会有几个普通User（注意，Admin是Admin，不包含在这个计数中），这里配置了1，也就是说我们只需要一个普通用户User1@org2.example.com 我们可以根据实际需要调整这个配置文件，增删Org Users等。

1.3生成公私钥和证书
我们配置好crypto-config.yaml文件后，就可以用cryptogen去读取该文件，并生成对应的公私钥和证书了：

cd examples/e2e_cli/
../../build/bin/cryptogen generate --config=./crypto-config.yaml
生成的文件都保存到crypto-config文件夹，我们可以进入该文件夹查看生成了哪些文件：

tree crypto-config
2.生成创世区块和Channel配置区块
2.1编译生成configtxgen
与前面1.1说到的类似，我们可以通过make命令生成configtxgen程序：

cd ~/go/src/github.com/hyperledger/fabric

make configtxgen
运行后的结果为：

build/bin/configtxgen 
CGO_CFLAGS=" " GOBIN=/home/studyzy/go/src/github.com/hyperledger/fabric/build/bin go install -tags "nopkcs11" -ldflags "-X github.com/hyperledger/fabric/common/configtx/tool/configtxgen/metadata.Version=1.0.0" github.com/hyperledger/fabric/common/configtx/tool/configtxgen 
Binary available as build/bin/configtxgen
2.2配置configtx.yaml
官方提供的examples/e2e_cli/configtx.yaml这个文件里面配置了由2个Org参与的Orderer共识配置TwoOrgsOrdererGenesis，以及由2个Org参与的Channel配置：TwoOrgsChannel。Orderer可以设置共识的算法是Solo还是Kafka，以及共识时区块大小，超时时间等，我们使用默认值即可，不用更改。而Peer节点的配置包含了MSP的配置，锚节点的配置。如果我们有更多的Org，或者有更多的Channel，那么就可以根据模板进行对应的修改。

2.3生成创世区块
配置修改好后，我们就用configtxgen 生成创世区块。并把这个区块保存到本地channel-artifacts文件夹中：

cd examples/e2e_cli/

../../build/bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
2.4生成Channel配置区块
../../build/bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
另外关于锚节点的更新，我们也需要使用这个程序来生成文件：

../../build/bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP

../../build/bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
最终，我们在channel-artifacts文件夹中，应该是能够看到4个文件。

channel-artifacts/
├── channel.tx
├── genesis.block
├── Org1MSPanchors.tx
└── Org2MSPanchors.tx

3.配置Fabric环境的docker-compose文件
前面对节点和用户的公私钥以及证书，还有创世区块都生成完毕，接下来我们就可以配置docker-compose的yaml文件，启动Fabric的Docker环境了。

3.1配置Orderer
Orderer的配置是在base/docker-compose-base.yaml里面，我们看看其中的内容：

复制代码
orderer.example.com: 
  container_name: orderer.example.com 
  image: hyperledger/fabric-orderer 
  environment: 
    - ORDERER_GENERAL_LOGLEVEL=debug 
    - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0 
    - ORDERER_GENERAL_GENESISMETHOD=file 
    - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block 
    - ORDERER_GENERAL_LOCALMSPID=OrdererMSP 
     - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp 
     # enabled TLS 
    - ORDERER_GENERAL_TLS_ENABLED=true 
    - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key 
    - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt 
    - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt] 
  working_dir: /opt/gopath/src/github.com/hyperledger/fabric 
  command: orderer 
  volumes: 
  - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block 
  - ../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp:/var/hyperledger/orderer/msp 
  - ../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/:/var/hyperledger/orderer/tls 
  ports: 
    - 7050:7050
复制代码
这里主要关心的是，ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block，而这个创世区块就是我们之前创建的创世区块，这里就是Host到Docker的映射：

  - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block

另外的配置主要是TL，Log等，最后暴露出服务端口7050。

3.2配置Peer
Peer的配置是在base/docker-compose-base.yaml和peer-base.yaml里面，我们摘取其中的peer0.org1看看其中的内容：

复制代码
peer-base: 
  image: hyperledger/fabric-peer 
  environment: 
    - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock 
    # the following setting starts chaincode containers on the same 
    # bridge network as the peers 
    # https://docs.docker.com/compose/networking/ 
    - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=e2ecli_default 
    #- CORE_LOGGING_LEVEL=ERROR 
    - CORE_LOGGING_LEVEL=DEBUG 
    - CORE_PEER_TLS_ENABLED=true 
    - CORE_PEER_GOSSIP_USELEADERELECTION=true 
    - CORE_PEER_GOSSIP_ORGLEADER=false 
    - CORE_PEER_PROFILE_ENABLED=true 
    - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt 
    - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key 
    - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt 
  working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer 
  command: peer node start

peer0.org1.example.com: 
  container_name: peer0.org1.example.com 
  extends: 
    file: peer-base.yaml 
    service: peer-base 
  environment: 
    - CORE_PEER_ID=peer0.org1.example.com 
    - CORE_PEER_ADDRESS=peer0.org1.example.com:7051 
    - CORE_PEER_CHAINCODELISTENADDRESS=peer0.org1.example.com:7052 
    - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org1.example.com:7051 
    - CORE_PEER_LOCALMSPID=Org1MSP 
  volumes: 
      - /var/run/:/host/var/run/ 
      - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/fabric/msp 
      - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls:/etc/hyperledger/fabric/tls 
  ports: 
    - 7051:7051 
     - 7052:7052 
    - 7053:7053
复制代码
在Peer的配置中，主要是给Peer分配好各种服务的地址，以及TLS和MSP信息。

3.3配置CLI
CLI在整个Fabric网络中扮演客户端的角色，我们在开发测试的时候可以用CLI来代替SDK，执行各种SDK能执行的操作。CLI会和Peer相连，把指令发送给对应的Peer执行。CLI的配置在docker-compose-cli.yaml中，我们看看其中的内容：

复制代码
cli: 
  container_name: cli 
  image: hyperledger/fabric-tools 
  tty: true 
  environment: 
    - GOPATH=/opt/gopath 
    - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock 
    - CORE_LOGGING_LEVEL=DEBUG 
    - CORE_PEER_ID=cli 
    - CORE_PEER_ADDRESS=peer0.org1.example.com:7051 
    - CORE_PEER_LOCALMSPID=Org1MSP 
    - CORE_PEER_TLS_ENABLED=true 
    - CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt 
    - CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key 
    - CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
    - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
  working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer 
  command: /bin/bash -c './scripts/script.sh ${CHANNEL_NAME}; sleep $TIMEOUT' 
  volumes: 
      - /var/run/:/host/var/run/ 
      - ../chaincode/go/:/opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/go 
      - ./crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ 
       - ./scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/ 
      - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts 
  depends_on: 
    - orderer.example.com 
    - peer0.org1.example.com 
    - peer1.org1.example.com 
    - peer0.org2.example.com 
    - peer1.org2.example.com 
复制代码
从这里我们可以看到，CLI启动的时候默认连接的是peer0.org1.example.com，并且启用了TLS。默认是以Admin@org1.example.com这个身份连接到Peer的。CLI启动的时候，会去执行./scripts/script.sh 脚本，这个脚本也就是fabric/examples/e2e_cli/scripts/script.sh 这个脚本，这个脚本完成了Fabric环境的初始化和ChainCode的安装及运行，也就是接下来要讲的步骤4和5.在文件映射配置上，我们注意到../chaincode/go/:/opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/go，也就是说我们要安装的ChainCode都是在fabric/examples/chaincode/go目录下，以后我们要开发自己的ChainCode，只需要把我们的代码复制到该目录即可。

【注意：请注释掉cli中command这一行，我们不需要CLI启动的时候自动执行脚本，我们在步骤4,5要一步步的手动执行！】

4.初始化Fabric环境
4.1启动Fabric环境的容器
我们将整个Fabric Docker环境的配置放在docker-compose-cli.yaml后，只需要使用以下命令即可：

docker-compose -f docker-compose-cli.yaml up -d
最后这个-d参数如果不加，那么当前终端就会一直附加在docker-compose上，而如果加上的话，那么docker容器就在后台运行。运行docker ps命令可以看启动的结果：

CONTAINER ID        IMAGE                        COMMAND             CREATED             STATUS              PORTS                                                                       NAMES
6f98f57714b5        hyperledger/fabric-tools     "/bin/bash"         8 seconds ago       Up 7 seconds                                                                                    cli
6e7b3fd0e803        hyperledger/fabric-peer      "peer node start"   11 seconds ago      Up 8 seconds        0.0.0.0:10051->7051/tcp, 0.0.0.0:10052->7052/tcp, 0.0.0.0:10053->7053/tcp   peer1.org2.example.com
9e67abfb982f        hyperledger/fabric-orderer   "orderer"           11 seconds ago      Up 8 seconds        0.0.0.0:7050->7050/tcp                                                      orderer.example.com
908d7fe2a4c7        hyperledger/fabric-peer      "peer node start"   11 seconds ago      Up 9 seconds        0.0.0.0:7051-7053->7051-7053/tcp                                            peer0.org1.example.com
6bb187ac10ec        hyperledger/fabric-peer      "peer node start"   11 seconds ago      Up 10 seconds       0.0.0.0:9051->7051/tcp, 0.0.0.0:9052->7052/tcp, 0.0.0.0:9053->7053/tcp      peer0.org2.example.com
150baa520ed0        hyperledger/fabric-peer      "peer node start"   12 seconds ago      Up 9 seconds        0.0.0.0:8051->7051/tcp, 0.0.0.0:8052->7052/tcp, 0.0.0.0:8053->7053/tcp      peer1.org1.example.com

可以看到1Orderer+4Peer+1CLI都启动了。

4.2创建Channel
现在我们要进入cli容器内部，在里面创建Channel。先用以下命令进入CLI内部Bash：

docker exec -it cli bash
创建Channel的命令是peer channel create，我们前面创建2.4创建Channel的配置区块时，指定了Channel的名字是mychannel，那么这里我们必须创建同样名字的Channel。

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls true --cafile $ORDERER_CA
执行该命令后，系统会提示：

2017-08-29 20:36:47.486 UTC [channelCmd] readBlock -> DEBU 020 Received block:0

系统会在cli内部的当前目录创建一个mychannel.block文件，这个文件非常重要，接下来其他节点要加入这个Channel就必须使用这个文件。

4.3各个Peer加入Channel
前面说过，我们CLI默认连接的是peer0.org1，那么我们要将这个Peer加入mychannel就很简单，只需要运行如下命令：

peer channel join -b mychannel.block
系统返回消息：

2017-08-29 20:40:27.053 UTC [channelCmd] executeJoin -> INFO 006 Peer joined the channel!

那么其他几个Peer又该怎么加入Channel呢？这里就需要修改CLI的环境变量，使其指向另外的Peer。比如我们要把peer1.org1加入mychannel，那么命令是：

CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer1.org1.example.com:7051

peer channel join -b mychannel.block
系统会返回成功加入Channel的消息。

同样的方法，将peer0.org2加入mychannel：

CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer channel join -b mychannel.block
最后把peer1.org2加入mychannel：

CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer1.org2.example.com:7051

peer channel join -b mychannel.block
4.4更新锚节点
关于AnchorPeer，我理解的不够深刻，经过我的测试，即使没有设置锚节点的情况下，整个Fabric网络仍然是能正常运行的。

对于Org1来说，peer0.org1是锚节点，我们需要连接上它并更新锚节点：

CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer0.org1.example.com:7051

peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx --tls true --cafile $ORDERER_CA
另外对于Org2，peer0.org2是锚节点，对应的更新代码是：

CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org2MSPanchors.tx --tls true --cafile $ORDERER_CA
5.链上代码的安装与运行
以上，整个Fabric网络和Channel都准备完毕，接下来我们来安装和运行ChainCode。这里仍然以最出名的Example02为例。这个例子实现了a，b两个账户，相互之间可以转账。

5.1Install ChainCode安装链上代码
链上代码的安装需要在各个相关的Peer上进行，对于我们现在这种Fabric网络，如果4个Peer都想对Example02进行操作，那么就需要安装4次。

仍然是保持在CLI的命令行下，我们先切换到peer0.org1这个节点：

CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer0.org1.example.com:7051
使用peer chaincode install命令可以安装指定的ChainCode并对其命名：

peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02
安装的过程其实就是对CLI中指定的代码进行编译打包，并把打包好的文件发送到Peer，等待接下来的实例化。

其他节点由于暂时还没使用到，我们可以先不安装，等到了步骤5.4的时候再安装。

5.2Instantiate ChainCode实例化链上代码
实例化链上代码主要是在Peer所在的机器上对前面安装好的链上代码进行包装，生成对应Channel的Docker镜像和Docker容器。并且在实例化时我们可以指定背书策略。我们运行以下命令完成实例化：

peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc -v 1.0 -c '{"Args":["init","a","100","b","200"]}' -P "OR      ('Org1MSP.member','Org2MSP.member')"
如果我们新开一个Ubuntu终端，去查看peer0.org1上的日志，那么就可以知道整个实例化的过程到底干了什么：

docker logs -f peer0.org1.example.com
主要几行重要的日志：

复制代码
2017-08-29 21:14:12.290 UTC [chaincode-platform] generateDockerfile -> DEBU 3fd 
FROM hyperledger/fabric-baseos:x86_64-0.3.1 
ADD binpackage.tar /usr/local/bin 
LABEL org.hyperledger.fabric.chaincode.id.name="mycc" \ 
       org.hyperledger.fabric.chaincode.id.version="1.0" \ 
      org.hyperledger.fabric.chaincode.type="GOLANG" \ 
      org.hyperledger.fabric.version="1.0.0" \ 
      org.hyperledger.fabric.base.version="0.3.1" 
ENV CORE_CHAINCODE_BUILDLEVEL=1.0.0 
ENV CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/peer.crt 
COPY peer.crt /etc/hyperledger/fabric/peer.crt 
2017-08-29 21:14:12.297 UTC [util] DockerBuild -> DEBU 3fe Attempting build with image hyperledger/fabric-ccenv:x86_64-1.0.0 
2017-08-29 21:14:48.907 UTC [dockercontroller] deployImage -> DEBU 3ff Created image: dev-peer0.org1.example.com-mycc-1.0 
2017-08-29 21:14:48.908 UTC [dockercontroller] Start -> DEBU 400 start-recreated image successfully 
2017-08-29 21:14:48.908 UTC [dockercontroller] createContainer -> DEBU 401 Create container: dev-peer0.org1.example.com-mycc-1.0
复制代码
接下来的日志就是各种初始化，验证，写账本之类的。总之完毕后，我们回到Ubuntu终端，使用docker ps可以看到有新的容器正在运行：

CONTAINER ID        IMAGE                                 COMMAND                  CREATED              STATUS              PORTS                                                                       NAMES
07791d4a99b7        dev-peer0.org1.example.com-mycc-1.0   "chaincode -peer.a..."   About a minute ago   Up About a minute                                                                               dev-peer0.org1.example.com-mycc-1.0
6f98f57714b5        hyperledger/fabric-tools              "/bin/bash"              About an hour ago    Up About an hour                                                                                cli
6e7b3fd0e803        hyperledger/fabric-peer               "peer node start"        About an hour ago    Up About an hour    0.0.0.0:10051->7051/tcp, 0.0.0.0:10052->7052/tcp, 0.0.0.0:10053->7053/tcp   peer1.org2.example.com
9e67abfb982f        hyperledger/fabric-orderer            "orderer"                About an hour ago    Up About an hour    0.0.0.0:7050->7050/tcp                                                      orderer.example.com
908d7fe2a4c7        hyperledger/fabric-peer               "peer node start"        About an hour ago    Up About an hour    0.0.0.0:7051-7053->7051-7053/tcp                                            peer0.org1.example.com
6bb187ac10ec        hyperledger/fabric-peer               "peer node start"        About an hour ago    Up About an hour    0.0.0.0:9051->7051/tcp, 0.0.0.0:9052->7052/tcp, 0.0.0.0:9053->7053/tcp      peer0.org2.example.com
150baa520ed0        hyperledger/fabric-peer               "peer node start"        About an hour ago    Up About an hour    0.0.0.0:8051->7051/tcp, 0.0.0.0:8052->7052/tcp, 0.0.0.0:8053->7053/tcp      peer1.org1.example.com

5.3在一个Peer上查询并发起交易
现在链上代码的实例也有了，并且在实例化的时候指定了a账户100，b账户200，我们可以试着调用ChainCode的查询代码，验证一下，在cli容器内执行：

peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
返回结果：Query Result: 100

接下来我们可以试着把a账户的10元转给b。对应的代码：

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n mycc -c '{"Args":["invoke","a","b","10"]}'
5.4在另一个节点上查询交易
前面的操作都是在org1下面做的，那么处于同一个区块链（同一个Channel下）的org2，是否会看org1的更改呢？我们试着给peer0.org2安装链上代码：

CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02
由于mycc已经在前面org1的时候实例化了，也就是说对应的区块已经生成了，所以在org2不能再次初始化。我们直接运行查询命令：

peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
这个时候我们发现运行该命令后要等很久（我这里花了40秒）才返回结果：

Query Result: 90

这是因为peer0.org2也需要生成Docker镜像，创建对应的容器，才能通过容器返回结果。我们回到Ubuntu终端，执行docker ps，可以看到又多了一个容器：

CONTAINER ID        IMAGE                                 COMMAND                  CREATED             STATUS              PORTS                                                                       NAMES
3e37aba50189        dev-peer0.org2.example.com-mycc-1.0   "chaincode -peer.a..."   2 minutes ago       Up 2 minutes                                                                                    dev-peer0.org2.example.com-mycc-1.0
07791d4a99b7        dev-peer0.org1.example.com-mycc-1.0   "chaincode -peer.a..."   21 minutes ago      Up 21 minutes                                                                                   dev-peer0.org1.example.com-mycc-1.0
6f98f57714b5        hyperledger/fabric-tools              "/bin/bash"              About an hour ago   Up About an hour                                                                                cli
6e7b3fd0e803        hyperledger/fabric-peer               "peer node start"        About an hour ago   Up About an hour    0.0.0.0:10051->7051/tcp, 0.0.0.0:10052->7052/tcp, 0.0.0.0:10053->7053/tcp   peer1.org2.example.com
9e67abfb982f        hyperledger/fabric-orderer            "orderer"                About an hour ago   Up About an hour    0.0.0.0:7050->7050/tcp                                                      orderer.example.com
908d7fe2a4c7        hyperledger/fabric-peer               "peer node start"        About an hour ago   Up About an hour    0.0.0.0:7051-7053->7051-7053/tcp                                            peer0.org1.example.com
6bb187ac10ec        hyperledger/fabric-peer               "peer node start"        About an hour ago   Up About an hour    0.0.0.0:9051->7051/tcp, 0.0.0.0:9052->7052/tcp, 0.0.0.0:9053->7053/tcp      peer0.org2.example.com
150baa520ed0        hyperledger/fabric-peer               "peer node start"        About an hour ago   Up About an hour    0.0.0.0:8051->7051/tcp, 0.0.0.0:8052->7052/tcp, 0.0.0.0:8053->7053/tcp      peer1.org1.example.com

总结
通过以上的分解，希望大家对Fabric环境的创建有了更深入的理解。我这里的示例仍然是官方的示例，并没有什么太新的东西。只要把这每一步搞清楚，那么接下来我们在生产环境创建更多的Org，创建大量的Channel，执行各种ChainCode都是如出一辙。