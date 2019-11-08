### 第八课操作题步骤——启动kafka排序服务操作

#### 节点启动目标

 实现一个基于 Kafka 提供排序服务的集群环境。集群包括2个 Org 组织，每个 Org 组织下各有2个 peer 节点（默认配置），4个 orderer节点，以及 Kafka 集群实现排序服务，由**4个 kafka 节点**，**3个 zookeeper 节点**组成。 

#### 1.配置文件修改（类似第四课作业）

需要配置crypto-config.yaml、configtx.yaml、docker-compose-base.yaml、docker-compose-cli.yaml四个文件

在进行配置文件修改前建议先将配置文件备份，以crypto-config.yaml为例，备份命令如下

```
# 进入到配置文件目录下
$ cd hyfa/fabric-samples/first-network/
# 复制crypto-config.yaml，并命名为crypto-config_backup.yaml
$ sudo cp crypto-config.yaml crypto-config_backup.yaml
```

##### 1.1 crypto-config.yaml修改

PeerOrgs中内容保持不变，修改OrderOrgs内容，修改如下

```
# ---------------------------------------------------------------------------
# "OrdererOrgs" - Definition of organizations managing orderer nodes
# ---------------------------------------------------------------------------
OrdererOrgs:
  # ---------------------------------------------------------------------------
  # Orderer
  # ---------------------------------------------------------------------------
  - Name: Orderer
    Domain: example.com
    # ---------------------------------------------------------------------------
    # "Specs" - See PeerOrgs below for complete description
    # ---------------------------------------------------------------------------
    Specs:
    # 四个orderer
      - Hostname: orderer
      - Hostname: orderer1
      - Hostname: orderer2
      - Hostname: orderer3
```

##### 1.2 configtx.yaml修改

备份原配置文件，并用vim编辑

```
$ sudo cp configtx.yaml configtx_backup.yaml
$ sudo vim configtx.yaml
```

主要修改内容包括：

 - Orderer.OrdererType 的值由默认的 **solo** 修改为 **kafka**
 - Orderer.Addresses 下添加另外的三个 orderer节点信息
 - Orderer.Kafka.Brokers 中添加 kafka 集群服务器的信息  

具体修改如下

```
################################################################################
#
#   SECTION: Orderer
#
#   - This section defines the values to encode into a config transaction or
#   genesis block for orderer related parameters
#
################################################################################
Orderer: &OrdererDefaults

    # Orderer Type: The orderer implementation to start
    # Available types are "solo" and "kafka"
    # 修改为kafka
    OrdererType: kafka

    Addresses:
    # 新增orderer节点
        - orderer.example.com:7050
        - orderer1.example.com:7050
        - orderer2.example.com:7050
        - orderer3.example.com:7050

    # Batch Timeout: The amount of time to wait before creating a batch
    BatchTimeout: 2s

    # Batch Size: Controls the number of messages batched into a block
    BatchSize:

        # Max Message Count: The maximum number of messages to permit in a batch
        MaxMessageCount: 10

        # Absolute Max Bytes: The absolute maximum number of bytes allowed for
        # the serialized messages in a batch.
        AbsoluteMaxBytes: 99 MB

        # Preferred Max Bytes: The preferred maximum number of bytes allowed for
        # the serialized messages in a batch. A message larger than the preferred
        # max bytes will result in a batch larger than preferred max bytes.
        PreferredMaxBytes: 512 KB

    Kafka:
        # Brokers: A list of Kafka brokers to which the orderer connects
        # NOTE: Use IP:port notation
        # 配置kafka集群服务器
        Brokers:
            - kafka0:9092
            - kafka1:9092
            - kafka2:9092
            - kafka3:9092

    # Organizations is the list of orgs which are defined as participants on
    # the orderer side of the network
    Organizations:
```

##### 1.3  docker-compose-base.yaml修改

备份原配置文件，并编辑

```
$ sudo cp base/docker-compose-base.yaml base/docker-compose-base_backup.yaml
$ sudo vim base/docker-compose-base.yaml
```

主要修改内容包括：

- 新增 zookeeper 节点信息、kafka 节点信息

- 修改 orderer 节点的信息

- peer 节点信息不变 

具体修改如下：

```
  zookeeper:
    image: hyperledger/fabric-zookeeper
    environment: 
      - ZOO_SERVERS=server.1=zookeeper1:2888:3888 server.2=zookeeper2:2888:3888 server.3=zookeeper3:2888:3888
    ports:
      - '2181'
      - '2888'
      - '3888'

  kafka:
    image: hyperledger/fabric-kafka
    environment:
      - KAFKA_LOG_RETENTION_MS=-1
      - KAFKA_MESSAGE_MAX_BYTES=103809024
      - KAFKA_REPLICA_FETCH_MAX_BYTES=103809024
      - KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false      
      - KAFKA_MIN_INSYNC_REPLICAS=2
      - KAFKA_DEFAULT_REPLICATION_FACTOR=3
      - KAFKA_ZOOKEEPER_CONNECT=zookeeper1:2181,zookeeper2:2181,zookeeper3:2181
    ports:
      - '9092'    

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
      # kafka      
      - CONFIGTX_ORDERER_ORDERERTYPE=kafka
      - CONFIGTX_ORDERER_KAFKA_BROKERS=[kafka0:9092,kafka1:9092,kafka2:9092,kafka3:9092]
      - ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s
      - ORDERER_KAFKA_RETRY_SHORTTOTAL=30s
      - ORDERER_KAFKA_VERBOSE=true    
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
      - '7050'

  orderer1.example.com:
    container_name: orderer1.example.com
    image: hyperledger/fabric-orderer
    environment:
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      # kafka      
      - CONFIGTX_ORDERER_ORDERERTYPE=kafka
      - CONFIGTX_ORDERER_KAFKA_BROKERS=[kafka0:9092,kafka1:9092,kafka2:9092,kafka3:9092]
      - ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s
      - ORDERER_KAFKA_RETRY_SHORTTOTAL=30s
      - ORDERER_KAFKA_VERBOSE=true
      # enabled TLS
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ../crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/msp:/var/hyperledger/orderer/msp
      - ../crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/:/var/hyperledger/orderer/tls
    ports:
      - '7050'

  orderer2.example.com:
    container_name: orderer2.example.com
    image: hyperledger/fabric-orderer
    environment:
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      # kafka      
      - CONFIGTX_ORDERER_ORDERERTYPE=kafka
      - CONFIGTX_ORDERER_KAFKA_BROKERS=[kafka0:9092,kafka1:9092,kafka2:9092,kafka3:9092]
      - ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s
      - ORDERER_KAFKA_RETRY_SHORTTOTAL=30s
      - ORDERER_KAFKA_VERBOSE=true
      # enabled TLS
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ../crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/msp:/var/hyperledger/orderer/msp
      - ../crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/:/var/hyperledger/orderer/tls
    ports:
      - '7050'

  orderer3.example.com:
    container_name: orderer3.example.com
    image: hyperledger/fabric-orderer
    environment:
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      # kafka      
      - CONFIGTX_ORDERER_ORDERERTYPE=kafka
      - CONFIGTX_ORDERER_KAFKA_BROKERS=[kafka0:9092,kafka1:9092,kafka2:9092,kafka3:9092]
      - ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s
      - ORDERER_KAFKA_RETRY_SHORTTOTAL=30s
      - ORDERER_KAFKA_VERBOSE=true
      # enabled TLS
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ../crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/msp:/var/hyperledger/orderer/msp
      - ../crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/:/var/hyperledger/orderer/tls
    ports:
      - '7050'
```

##### 1.4 docker-compose-cli.yaml修改

同样我们备份原配置文件，并vim编辑

```
$ sudo cp docker-compose-cli.yaml docker-compose-cli_backup.yaml
$ sudo vim docker-compose-cli.yaml
```

主要修改内容包括：

- 添加三个 zookeeper 节点

- 添加四个 kafka 节点

- 添加相应的 orderer 节点的信息 

```
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

version: '2'

volumes:
  orderer.example.com:
  orderer1.example.com:
  orderer2.example.com:
  orderer3.example.com:
  peer0.org1.example.com:
  peer1.org1.example.com:
  peer0.org2.example.com:
  peer1.org2.example.com:

networks:
  byfn:

services:

  zookeeper1:
    container_name: zookeeper1
    extends:
      file:   base/docker-compose-base.yaml
      service: zookeeper
    environment:
      - ZOO_MY_ID=1
    networks:
      - byfn

  zookeeper2:
    container_name: zookeeper2
    extends:
      file:   base/docker-compose-base.yaml
      service: zookeeper
    environment:
      - ZOO_MY_ID=2
    networks:
      - byfn

  zookeeper3:
    container_name: zookeeper3
    extends:
      file:   base/docker-compose-base.yaml
      service: zookeeper
    environment:
      - ZOO_MY_ID=3
    networks:
      - byfn

  kafka0:
    container_name: kafka0
    extends:
      file:   base/docker-compose-base.yaml
      service: kafka
    environment:
      - KAFKA_BROKER_ID=0      
    depends_on:
      - zookeeper1
      - zookeeper2
      - zookeeper3
    networks:
      - byfn

  kafka1:
    container_name: kafka1
    extends:
      file:   base/docker-compose-base.yaml
      service: kafka
    environment:
      - KAFKA_BROKER_ID=1
    depends_on:
      - zookeeper1
      - zookeeper2
      - zookeeper3
    networks:
      - byfn

  kafka2:
    container_name: kafka2
    extends:
      file:   base/docker-compose-base.yaml
      service: kafka
    environment:
      - KAFKA_BROKER_ID=2
    depends_on:
      - zookeeper1
      - zookeeper2
      - zookeeper3
    networks:
      - byfn

  kafka3:
    container_name: kafka3
    extends:
      file:   base/docker-compose-base.yaml
      service: kafka
    environment:
      - KAFKA_BROKER_ID=3
    depends_on:
      - zookeeper1
      - zookeeper2
      - zookeeper3
    networks:
      - byfn

  orderer.example.com:
    extends:
      file:   base/docker-compose-base.yaml
      service: orderer.example.com
    container_name: orderer.example.com
    depends_on:
      - kafka0
      - kafka1
      - kafka2
      - kafka3
    networks:
      - byfn

  orderer1.example.com:
    extends:
      file:   base/docker-compose-base.yaml
      service: orderer1.example.com
    container_name: orderer1.example.com
    depends_on:
      - kafka0
      - kafka1
      - kafka2
      - kafka3
    networks:
      - byfn

  orderer2.example.com:
    extends:
      file:   base/docker-compose-base.yaml
      service: orderer2.example.com
    container_name: orderer2.example.com
    depends_on:
      - kafka0
      - kafka1
      - kafka2
      - kafka3
    networks:
      - byfn

  orderer3.example.com:
    extends:
      file:   base/docker-compose-base.yaml
      service: orderer3.example.com
    container_name: orderer3.example.com   
    depends_on:
      - kafka0
      - kafka1
      - kafka2
      - kafka3
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
    depends_on:
      - orderer.example.com
      - orderer1.example.com
      - orderer2.example.com
      - orderer3.example.com
      - peer0.org1.example.com
      - peer1.org1.example.com
      - peer0.org2.example.com
      - peer1.org2.example.com
    networks:
      - byfn
```

#### 2. 启动网络

同理第三节课中启动步骤

使用 byfn.sh **生成组织结构及身份证书及所需的各项配置文件**，命令如下

```
$ sudo ./byfn.sh generate
```

**启动网络**

```
$ sudo docker-compose -f docker-compose-cli.yaml up
```

最后可通过docker ps 查看启动的3个zookeeper节点、4个kafka节点和4个orderer节点
