# Course 04

## contents

### 4逐步启动网络

前面的课程中，我们通过自动化脚本 byfn.sh 可以自动帮我们创建网络环境运行时所需的所有内容，但在一些特定情况之下，我们根据不同的需求需要自定义一些设置。因此我们来逐步启动网络了解其启动原理和细节。

#### 4.1 生成组织结构

生成过程依赖`crypto-config.yaml` 配置文件，该配置文件路径 ：`~/fabric-samples/first-network/crypto-config.yaml`  （~代表fabric-samples文件夹所在的父目录）

crypto-config.yaml 配置文件包含如下内容：

```
OrdererOrgs:
  - Name: Orderer    # Orderer的名称
    Domain: example.com    # 域名
    Specs:
      - Hostname: orderer    # hostname + Domain的值组成Orderer节点的完整域名

PeerOrgs:
  - Name: Org1
    Domain: org1.example.com
    EnableNodeOUs: true        # 在msp下生成config.yaml文件
    Template:                  # 组织中的节点数目
      Count: 2
    Users:                     #组织中的用户数目
      Count: 1

  - Name: Org2
    Domain: org2.example.com
    EnableNodeOUs: true
    Template:
      Count: 2
    Users:
      Count: 1
```

该配置文件指定了 OrdererOrgs 及 PeerOrgs 两个组织信息。在 PeerOrgs 配置信息中指定了 Org1 与 Org2 两个组织。每个组织使用 **Template** 属性下的 Count 指定了两个节点， **Users** 属性下的 Count 指定了一个用户。

组织信息中还包含组织名称、域名、节点数量、及新增的用户数量等相关信息。

Peer 节点的域名组成为 peer + 起始数字0 + Domain，如 Org1 中的两个 peer 节点的完整域名为：

**peer0.org1.example.com**，**peer1.org1.example.com**

接下来依据Fabric提供的**cryptogen**模块和crypto-config.yaml配置文件来生成组织关系和身份证书

① 进入`fabric-samples/first-network` 目录：

```
$ cd hyfa/fabric-samples/first-network/
```

② 使用cryptogen工具生成指定拓扑结构的**组织关系和身份证书**，具体细节由crypto-config.yaml配置文件来指定

```
$ sudo ../bin/cryptogen generate --config=./crypto-config.yaml
```

执行后输出如下

```
org1.example.com
org2.example.com
```

#### 4.2 证书和密钥

身份证书将被输出到 ~/fabric-samples/first-network路径下的 `crypto-config` 的目录中

该目录下有两个子目录：**ordererOrganizations** 和 **peerOrganizations**

**ordererOrganizations** 子目录下包括构成 Orderer 组织(1个 Orderer 节点)的身份信息

**peerOrganizations** 子目录下为所有的 Peer 节点组织(2个组织，4个节点)的相关身份信息. 其中最关键的是 MSP 目录, 代表了实体的身份信息

在生成的目录结构中最关键的是各个资源下的 msp 目录内容，存储了生成的代表 MSP 实体身份的各种证书文件，一般包括：

- admincerts ：管理员的身份证书文件
- cacerts ：信任的根证书文件
- keystore ：节点的签名私钥文件
- signcerts ：节点的签名身份证书文件
- tlscacerts:：TLS 连接用的证书
- intermediatecerts （可选）：信任的中间证书
- crls （可选）：证书撤销列表
- config.yaml （可选）：记录OrganizationalUnitldentifiers 信息，包括根证书位置和ID信息

这些身份文件随后可以分发到对应的Orderer 节点和Peer 节点上，并放到对应的MSP路径下，用于签名验证使用。

#### 4.3 创世区块和通道的建立

创建服务启动初始区块及应用通道交易配置文件定义在名为 **configtx.yaml** 文件中，关键的配置信息注释如下

```
---
Organizations:
    - &OrdererOrg
        Name: OrdererOrg
        ID: OrdererMSP
        MSPDir: crypto-config/ordererOrganizations/example.com/msp

    - &Org1
        Name: Org1MSP
        ID: Org1MSP
        MSPDir: crypto-config/peerOrganizations/org1.example.com/msp

        AnchorPeers:
            - Host: peer0.org1.example.com 
              Port: 7051

    - &Org2
        Name: Org2MSP
        ID: Org2MSP
        MSPDir: crypto-config/peerOrganizations/org2.example.com/msp

        AnchorPeers:
            - Host: peer0.org2.example.com   #新增组织Org3时 注意域名的命名
              Port: 7051

Capabilities:
    Global: &ChannelCapabilities
        V1_1: true

    Orderer: &OrdererCapabilities
        V1_1: true

    Application: &ApplicationCapabilities
        V1_2: true

Application: &ApplicationDefaults
    Organizations:

Orderer: &OrdererDefaults
    OrdererType: solo   //排序服务 solo/kafka
    Addresses:   //指定了 Orderer 节点的服务地址与端口号
        - orderer.example.com:7050

    BatchTimeout: 2s
    BatchSize:   //指定了批处理大小，如最大交易数量，最大字节数及建议字节数
        MaxMessageCount: 10
        AbsoluteMaxBytes: 99 MB
        PreferredMaxBytes: 512 KB

    Kafka:
        Brokers:
            - 127.0.0.1:9092

    Organizations:

Profiles:  //指定了两个模板：TwoOrgsOrdererGenesis 与 TwoOrgsChannel 。
    TwoOrgsOrdererGenesis:  //模板一
        Capabilities: // 指定通道的权限信息。
            <<: *ChannelCapabilities
        Orderer:  //指定了Orderer服务的信息（OrdererOrg）及权限信息。
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
            Capabilities:
                <<: *OrdererCapabilities
        Consortiums:   //定义了联盟组成成员（Org1&Org2）
            SampleConsortium:
                Organizations:
                    - *Org1
                    - *Org2
    TwoOrgsChannel: //模板二
        Consortium: SampleConsortium   //指定了联盟信息
        Application:   //指定了组织及权限信息。
            <<: *ApplicationDefaults
            Organizations:
                - *Org1
                - *Org2
            Capabilities:
                <<: *ApplicationCapabilities
```

​        该配置文件中由 **Organizations** 定义了三个成员 Orderer Org、Org1、Org2，并且设置每个成员的MSP 目录的位置。而且为每个 PeerOrg 指定了相应的锚节点（Org1 组织中`peer0.org1.example.com`与 Org2 组织中`peer0.org2.example.com`）

Orderer部分指定的Orderer节点的信息

1. **OrdererType** 指定了共识排序服务的实现方式，有两种选择（solo 及 Kafka）。
2. **Addresses** 指定了 Orderer 节点的服务地址与端口号。
3. **BatchSize** 指定了批处理大小，如最大交易数量，最大字节数及建议字节数。

**Profiles** 部分指定了两个模板：TwoOrgsOrdererGenesis 与 TwoOrgsChannel 。

1. **TwoOrgsOrdererGenesis** 模板用来生成Orderer服务的初始区块文件，该模板由三部分组成：
   1.1 Capabilities 指定通道的权限信息。

   1.2 Orderer 指定了Orderer服务的信息（OrdererOrg）及权限信息。

   1.3 Consortiums 定义了联盟组成成员（Org1&Org2）。

2. **TwoOrgsChannel** 模板用来生成应用通道交易配置文件。由两部分组成：

   2.1 Consortium 指定了联盟信息。

   2.2 Application 指定了组织及权限信息。

④ 使用 `configtx.yaml` 文件中定义的 `TwoOrgsOrdererGenesis` 模板,，生成 Orderer 服务系统通道的初始区块文件。执行如下命令

```
$ sudo ../bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
```

⑤ 接着创建通道，后面重复使用同一个通道名，我们这里指定通道名称的环境变量，方便复用

```
$ export CHANNEL_NAME=mychannel
```

⑥ 使用 `configtx.yaml` 配置文件中的 `TwoOrgsChannel` 模板, 来生成新建通道的配置交易文件, 执行如下

```
$ sudo ../bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
```

（若有警告 可以先无视，不影响后续操作）

⑦ 生成锚节点更新配置文件，同样基于 `configtx.yaml` 配置文件中的 TwoOrgsChannel 模板，为每个组织分别生成锚节点更新配置，且注意指定对应的组织名称。依次执行如下命令

```
$ sudo ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

$ sudo ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
```

上述所有命令执行完成后，channel-artifacts目录下会有4个被创建的文件，目录结构如下

```
channel-artifacts/
├── channel.tx
├── genesis.block
├── Org1MSPanchors.tx
└── Org2MSPanchors.tx
```

#### 4.4 网络服务配置

通过一条命令来方便的启动 Hyperledger Fabric 网络中所有节点。

⑧ 如下命令所示

```
$ sudo docker-compose -f docker-compose-cli.yaml up -d
```

启动命令中主要涉及到三个配置文件，

 - docker-compose 工具的示例配置文件：**docker-compose-cli.yaml**
 - 指定Orderer和Peers节点的相关信息：**docker-compose-base.yaml**

- 设置了所有 peer 容器的基本的共同信息：**peer-base.yaml**

##### 4.4.1 docker-compose-cli.yaml

Hyperledger Fabric 采用了容器技术，所以需要一个简化的方式来集中化管理这这些节点容器，我们使用 docker-compose 这个工具个来实现一步到位的节点容器管理。

fabric-samples/first-network 目录下，文件名称为： docker-compose-cli.yaml 为docker-compose的配置文件

```
version: '2'

volumes:
  orderer.example.com:
  peer0.org1.example.com:
  peer1.org1.example.com:
  peer0.org2.example.com:
  peer1.org2.example.com:

networks:
  byfn:

services:  # 六个容器，一个 Orderer，属于两个 Orgs 组织的四个 Peer，还有一个 CLI
           # 若增加org3（两个节点） 则再新增两个容器

  orderer.example.com:
    extends:
      file:   base/docker-compose-base.yaml
      service: orderer.example.com
    container_name: orderer.example.com
    networks:
      - byfn

  peer0.org1.example.com:
    container_name: peer0.org1.example.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer0.org1.example.com
    networks:
      - byfn

  peer1.org1.example.com:
    container_name: peer1.org1.example.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer1.org1.example.com
    networks:
      - byfn

  peer0.org2.example.com:
    container_name: peer0.org2.example.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer0.org2.example.com
    networks:
      - byfn

  peer1.org2.example.com:
    container_name: peer1.org2.example.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer1.org2.example.com
    networks:
      - byfn

  cli:
    container_name: cli
    image: hyperledger/fabric-tools:$IMAGE_TAG
    tty: true
    stdin_open: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      #- CORE_LOGGING_LEVEL=DEBUG
      - CORE_LOGGING_LEVEL=INFO
      - CORE_PEER_ID=cli
      - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
      - CORE_PEER_LOCALMSPID=Org1MSP
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
      - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: /bin/bash
    volumes:
        - /var/run/:/host/var/run/
        - ./../chaincode/:/opt/gopath/src/github.com/chaincode
        - ./crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
        - ./scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/
        - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
    depends_on:  #指定了所依赖的相关容器
      - orderer.example.com
      - peer0.org1.example.com
      - peer1.org1.example.com
      - peer0.org2.example.com
      - peer1.org2.example.com
    networks:
      - byfn
```

##### 4.4.2 docker-compose-base.yaml

```
version: '2'

services:

  orderer.example.com:
    container_name: orderer.example.com
    image: hyperledger/fabric-orderer:$IMAGE_TAG
    environment:
      - ORDERER_GENERAL_LOGLEVEL=INFO
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
    - orderer.example.com:/var/hyperledger/production/orderer
    ports:
      - 7050:7050

  peer0.org1.example.com:
    container_name: peer0.org1.example.com
    extends:
      file: peer-base.yaml
      service: peer-base
    environment:
      - CORE_PEER_ID=peer0.org1.example.com
      - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer1.org1.example.com:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org1.example.com:7051
      - CORE_PEER_LOCALMSPID=Org1MSP
    volumes:
        - /var/run/:/host/var/run/
        - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/fabric/msp
        - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls:/etc/hyperledger/fabric/tls
        - peer0.org1.example.com:/var/hyperledger/production
    ports:
      - 7051:7051
      - 7053:7053

  peer1.org1.example.com:
    container_name: peer1.org1.example.com
    extends:
      file: peer-base.yaml
      service: peer-base
    environment:
      - CORE_PEER_ID=peer1.org1.example.com
      - CORE_PEER_ADDRESS=peer1.org1.example.com:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer1.org1.example.com:7051
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.org1.example.com:7051
      - CORE_PEER_LOCALMSPID=Org1MSP
    volumes:
        - /var/run/:/host/var/run/
        - ../crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp:/etc/hyperledger/fabric/msp
        - ../crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls:/etc/hyperledger/fabric/tls
        - peer1.org1.example.com:/var/hyperledger/production

    ports:
      - 8051:7051
      - 8053:7053

  peer0.org2.example.com:
    container_name: peer0.org2.example.com
    extends:
      file: peer-base.yaml
      service: peer-base
    environment:
      - CORE_PEER_ID=peer0.org2.example.com
      - CORE_PEER_ADDRESS=peer0.org2.example.com:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org2.example.com:7051
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer1.org2.example.com:7051
      - CORE_PEER_LOCALMSPID=Org2MSP
    volumes:
        - /var/run/:/host/var/run/
        - ../crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp:/etc/hyperledger/fabric/msp
        - ../crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls:/etc/hyperledger/fabric/tls
        - peer0.org2.example.com:/var/hyperledger/production
    ports:
      - 9051:7051
      - 9053:7053

  peer1.org2.example.com:
    container_name: peer1.org2.example.com
    extends:
      file: peer-base.yaml
      service: peer-base
    environment:
      - CORE_PEER_ID=peer1.org2.example.com
      - CORE_PEER_ADDRESS=peer1.org2.example.com:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer1.org2.example.com:7051
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.org2.example.com:7051
      - CORE_PEER_LOCALMSPID=Org2MSP
    volumes:
        - /var/run/:/host/var/run/
        - ../crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/msp:/etc/hyperledger/fabric/msp
        - ../crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls:/etc/hyperledger/fabric/tls
        - peer1.org2.example.com:/var/hyperledger/production
    ports:
      - 10051:7051
      - 10053:7053
```

Orderer 设置如下信息：

- **environment：**指定日志级别、监听地址、生成初始区块的提供方式、初始区块配置文件路径、本地 MSPID 及对应的目录、开启 TLS 验证及对应的证书、私钥信息等诸多重要信息。
- **working_dir：**进入容器后的默认工作目录
- **volumes：**指定系统中的初始区块配置文件、MSP、TLS目录映射到容器中的指定路径下。
- **ports：** 指定当前节点的监听端口。

各 Peers 设置了如下信息：

- **extends：**基本信息来源于哪个文件。

- **environment：**指定了容器的的 ID、监听地址及端口号、本地 MSPID。
- **volumes：**将系统的 msp 及 tls 目录映射到容器中的指定路径下。
- **ports：** 指定当前节点的监听端口。（**新增组织节点时 注意不要占用已用端口**）

##### 4.4.3 peer-base.yaml

该配置文件设置了所有 peer 容器的基本的共同信息，日志级别，是否开启 TLS 验证，是否采用 Leader 选举， 是否将当前节点设为 Leader， TLS 证书、私钥、根证书的路径、容器的默认工作路径、容器启动命令。
