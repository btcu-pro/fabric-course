# course 07 深入MSP成员管理与Fabric CA服务实现
[TOC]
## MSP的概念
###MSP定义及作用
在 Hyperledger Fabric 中，各个网络参与者之间的通信安全依赖于 PKI 标准来实现，并确保在区块链上发布的消息得到相应的认证。

PKI（Public Key Infrastructure）：公钥基础结构。由向各方（如服务的用户，服务提供商）发布数字证书的证书颁发机构组成，然后他们使用它们在与其环境交换的消息中对自己进行身份验证。

PKI 有四个关键要素：

数字证书：包含与证书持有者相关的一组属性的文档。最常见的证书类型是符合X.509标准的证书，允许在其结构中编码一方的识别细节。
公钥和私钥：身份验证和消息完整性是安全通信中的重要概念。身份验证要求确保交换消息的各方创建特定消息的身份。对于具有“完整性”的消息意味着在其传输期间不能被修改。
证书颁发机构：证书颁发机构向不同的参与者分发证书，这些证书由CA进行数字签名。CA是为组织的参与者提供可验证的数字身份的基础。
证书撤销列表：某种原因而被撤销的证书的引用列表。
PKI 只是一个体系结构，负责生成及颁发；在 Hyperledger Fabric 中的默认 MSP 实际上是使用符合 X.509 标准的证书作为身份，采用传统的公钥基础结构（PKI）分层模型来实现。

MSP（Membership Service Provider）：成员服务提供商，是 Hyperledger Fabric 1.0版本开始抽象出来的一个模块化组件。用于定义身份验证，进行身份验证和允许访问网络的规则。更确切地说，MSP 是 Hyperledger Fabric 对网络中的组成成员进行身份管理与验证的模块组件。

具体作用如下：

MSP 管理用户 ID。
验证想要加入网络的节点：每一个想加入网络中的节点必须提供其有效且合法的 MSP 信息。
为客户发起的交易提供凭证：在各节点（Client、Peer、Orderer）之间进行数据传输时，需要验证各节点的签名。
MSP 在 Hyperledger Fabric 中的分类：

**网络MSP**：对整个 Hyperledger Fabric 网络中的成员进行管理；定义参与组织的 MSP ，以及组织成员中的哪些成员被授权执行管理任务（如创建通道）

**通道MSP**:对一个通道中的组织成员进行管理。通道在特定的一组组织之间提供私有通信。在该通道的 MSP 环境中通道策略定义了谁有权限参与通道上的某些行为（如添加组织或实例化链码）。

**Peer MSP**：本地 MSP 在每个 Peer 的文件系统上定义，并且每个 Peer 都有一个单独的 MSP 实例。执行与通道 MSP 完全相同的功能，其限制是它仅适用于定义它的 Peer。

**Orderer MSP**：与 Peer MSP 相同，Orderer 本地 MSP 也在其节点的文件系统上定义，仅适用于该节点。

**User MSP**： 每一个组织都可以拥有多个不同的用户，都在其 Organizations 节点的文件系统上定义，仅适用该组织（包括该组织下的所有 Peer 节点）。

##  MSP 的组成结构
MSP的逻辑结构如下所示（与实际的物理结构会有所不同）：



如上图所示，MSP有九个元素。其中MSP名称是根文件夹名称，每个子文件夹代表MSP配置的不同元素：

**根CA（Root CAs）**：文件夹中包含根CA（CA：Certificate Authorities）的自签名 X.509 证书列表。用于自签名及给中间 CA 证书签名。

**中间CA（ICA）**：包含由根据 CA 颁发的证书列表。

**组织单位（OUs）**：这些单位列在 $FABRIC_CFG_PATH/msp/config.yaml 文件中，包含一个组织单位列表，其成员被视为该MSP所代表的组织的一部分。

**管理员（B）**：此文件夹包含一个标识列表，用于定义具有此组织管理员角色的角色。对于标准MSP 类型，此列表中应该有一个或多个 X.509 证书。

需要注意，仅仅一个具有管理员的角色，并不意味着他们可以管理特定的资源，给定标识在管理系统方面的实际功能由管理系统资源的策略决定。

**撤销证书（ReCA）**：保存已被撤销参与者身份的信息。

**签名证书（SCA）**：背书节点在交易提案响应中的签名证书。此文件夹对于本地 MSP 是必需的，并且该节点必须只有一个 X.509 证书。

**私钥（KeyStore）**：此文件夹是为 Peer 或 Orderer 节点（或客户端的本地MSP）的本地MSP定义的，并包含节点的签名密钥。此密钥以加密方式匹配 SCA 文件夹中包含的签名证书，并用于签署数据（如签署交易提议响应，作为认可阶段的一部分）。此文件夹对于本地MSP是必需的，并且必须只包含一个私钥。

**TLS根CA（TLS RCA）**：包含组织信任的用于 TLS 通信的根 CA 的自签名 X.509 证书列表。此文件夹中必须至少有一个 TLS 根 CA X.509 证书。

**TLS中间CA（TLS ICA）**：保存由 TLS 根 CA 颁发的中间证书列表。

### MSP应用
要想初始化一个MSP实例，每一个peer节点和orderer节点都需要在本地指定其配置并启动。

首先， 为了方便地在网络中引用MSP，每个MSP都需要一个特定的名字（如 OrdererMSP、Org1MSP 或 Org2MSP.domain.com）。此名字被称之为 MSP 标识符或 MSP ID。对于每个 MSP 实例来说，MSP 标识符都必须独一无二。

在系统起始阶段，需要指定在网络中出现的所有 MSP 的验证参数，且这些参数需要在系统通道的创世区块中指定。MSP的验证参数包括MSP标识符、信任源证书、中间 CA 和管理员的证书，以及 OU 说明和 CLR。系统的创世区块会在 orderer 节点设置阶段被提供给它们，且允许它们批准创建通道的请求。如果创世区块包含两个有相同标识符的 MSP，那么 orderer 节点将拒绝系统创世区块，导致网络引导程序执行失败。

要想生成 X.509 证书以满足 MSP 配置，应用程序可以有多种方式实现：

使用Openssl。在此需要注意：在 Hyperledger Fabric 中，不支持包括RSA密钥在内的证书。
使用 cryptogen 工具，其操作方法参见第三章 3.1 生成组织结构与身份证书 一节 。
Hyperledger Fabric CA 也可用于生成配置 MSP 所需的密钥及证书。详见下节内容。
在节点的配置文件中（对 peer 节点而言配置文件是 core.yaml 文件，对 orderer 节点而言则是orderer.yaml文件。在实际开发中可自定义配置文件名称），我们需要指定到 mspconfig 文件夹的路径，以及节点的 MSP 的 MSP 标识符。节点的 MSP 的 MSP 标识符则会作为参数 localMspId 和 LocalMSPID 的值分别提供给 peer 节点和 orderer 节点。

运行环境可以通过为 peer 使用 CORE 前缀（如 CORE_PEER_LOCALMSPID）及为 orderer 使用 ORDERER 前缀（例如 ORDERER_GENERAL_LOCALMSPID）对以上变量进行重写。如在 fabric-samples 中提供的示例配置文件 docker-compose-base.yaml：
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
      # 指定本地 MSP ID
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
    # MSP 映射信息
    - ../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp:/var/hyperledger/orderer/msp
    # TLS 映射信息
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
      # 指定本地 MSP ID
      - CORE_PEER_LOCALMSPID=Org1MSP

    volumes:
        - /var/run/:/host/var/run/
        # MSP 映射信息
        - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/fabric/msp
         # TLS 映射信息
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
      # 指定本地 MSP ID
      - CORE_PEER_LOCALMSPID=Org1MSP
    volumes:
        - /var/run/:/host/var/run/
         # MSP 映射信息
        - ../crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp:/etc/hyperledger/fabric/msp
         # TLS 映射信息
        - ../crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls:/etc/hyperledger/fabric/tls
        - peer1.org1.example.com:/var/hyperledger/production
    ports:
      - 8051:7051
      - 8053:7053

  ......
```

## Fabric CA的概念

###Fabric CA简介
Hyperledger Fabric CA 是 Hyperledger Fabric 的证书颁发机构（CA），是超级账本 Hyperledger Fabric 内一个可选的 MemberService 组件，对网络内各个实体的身份证书进行管理，主要实现：

负责 Fabric 网络内所有实体（Identity）身份的注册。
负责对数字证书的签发，包括 ECerts（身份证书）、TCerts（交易证书）。
证书的续签或吊销。
Fabric CA 在 Hyperledger Fabric 网络中的作用如下图所示：

访问 Fabric CA 服务器可以通过 Hyperledger Fabric CA 客户端或通过其中一个 Fabric SDK 来实现，与 Hyperledger Fabric CA 服务器的所有通信都是通过 REST API 进行。

Hyperledger Fabric CA 客户端或 SDK 可以连接到 Hyperledger Fabric CA 服务器集群，集群由 HA Proxy 等实现负载均衡。服务器可能包含多个CA，每个CA都是根CA或中间CA，每个中间CA都有一个父CA。

Hyperledger Fabric CA 的身份信息保存在数据库或LDAP中。目前 Fabric CA 支持的数据库有 MySQL、PostgreSQL、SQLite；默认使用 SQLite 数据库。如果配置了 LDAP，则身份信息将保留在 LDAP 而不是数据库中。

关于 Hyperledger Fabric CA 的更多详细信息，请 点击此处

### Fabric CA 安装
#### 环境要求
安装 Go1.9 或以上版本并设置 GOPATH 环境变量

安装 libtool 与 libltdl-dev 依赖包

$ sudo apt update
$ sudo apt install libtool libltdl-dev
如果没有安装 libtool libltdl-dev 依赖，会在安装 Fabric CA 时产生错误

#### 安装服务端与客户端
**方式一**：

安装服务端与客户端二进制命令到 $GOPATH/bin 目录下
```
$ go get -u github.com/hyperledger/fabric-ca/cmd/...
```
命令执行完成后，会自动在 $GOPATH/bin 目录下产生两个可执行文件：
```
fabric-ca-client
fabric-ca-server
```
设置环境变量，以便于在任何路径下都可以直接使用两个命令：

```
export PATH=$PATH:$GOPATH/bin
```
**方式二**：

除了如上方式外，还可以在 fabric-ca 目录下生成 fabric-ca-client、fabric-ca-server 两个可执行文件，方法如下：

切换至源码目录下：
```
$ cd $GOPATH/src/github.com/hyperledger/fabric-ca/
```
使用make命令编译：
```
$ make fabric-ca-server
$ make fabric-ca-client
```
自动在当前的 fabric-ca 目录下生成 bin 目录, 目录中包含 fabric-ca-client 与 fabric-ca-server 两个可执行文件。

设置环境变量：
```
$ export PATH=$GOPATH/src/github.com/hyperledger/fabric-ca/bin:$PATH
```

## 启动Fabric CA
###初始化
Fabric CA 服务器的主目录确定如下：

如果设置了 -home 命令行选项，则使用其值
否则，如果 FABRIC_CA_SERVER_HOME 设置了环境变量，则使用其值
否则，如果 FABRIC_CA_HOME 设置了环境变量，则使用其值
否则，如果 CA_CFG_PATH 设置了环境变量，则使用其值
否则，使用当前工作目录作为服务器端的主目录.
现在我们使用一个当前所在的目录作为服务器端的主目录。返回至用户的HOME目录下，创建一个 fabric-ca 目录并进入该目录
```
$ cd ~
$ mkdir fabric-ca
$ cd fabric-ca
```
创建该目录的目的是作为 Fabric CA 服务器的主目录。默认服务器主目录为 “./” 。

初始化 Fabric CA
```
$ fabric-ca-server init -b admin:pass
```
在初始化时 -b 选项是必需的，用于指定注册用户的用户名与密码。

命令执行后会自动生成配置文件到至当前目录：
```
fabric-ca-server-config.yaml： 默认配置文件
ca-cert.pem： PEM 格式的 CA 证书文件, 自签名
fabric-ca-server.db： 存放数据的 sqlite3 数据库
msp/keystore/： 路径下存放个人身份的私钥文件(_sk文件)，对应签名证书
```
#### 快速启动
快速启动并初始化一个 fabric-ca-server 服务
```
$ fabric-ca-server start -b admin:pass
```
-b： 提供注册用户的名称与密码, 如果没有使用 LDAP，这个选项为必需。默认的配置文件的名称为 fabric-ca-server-config.yaml

如果之前没有执行初始化命令, 则启动过程中会自动进行初始化操作. 即从主配置目录搜索相关证书和配置文件, 如果不存在则会自动生成

#### 配置数据库
Fabric CA 默认数据库为 SQLite，默认数据库文件 fabric-ca-server.db 位于 Fabric CA 服务器的主目录中。SQLite 是一个嵌入式的小型的数据系统，但在一些特定的情况下，我们需要集群来支持，所以Fabric CA 也设计了支持其它的数据库系统（目前只支持 MySQL、PostgreSQL 两种）。Fabric CA 在集群设置中支持以下数据库版本：

PostgreSQL：9.5.5 或更高版本
MySQL：5.7 或更高版本
下面我们来看如何配置来实现对不同数据库的支持。

##### 配置 PostgreSQL
如果使用 PostgreSQL 数据库，则需要在 Fabric CA 服务器端的配置文件进行如下设置：
```
db:
  type: postgres
  datasource: host=localhost port=5432 user=Username password=Password dbname=fabric_ca sslmode=verify-full
如果要使用 TLS，则必须指定 Fabric CA 服务器配置文件中的 db.tls 部分。如果在 PostgreSQL 服务器上启用了 SSL 客户端身份验证，则还必须在 db.tls.client 部分中指定客户端证书和密钥文件。如下所示：

db:
  ...
  tls:
      enabled: true
      certfiles:
        - db-server-cert.pem
      client:
            certfile: db-client-cert.pem
            keyfile: db-client-key.pem
certfiles：PEM 编码的受信任根证书文件列表。

certfile和keyfile：Fabric CA 服务器用于与 PostgreSQL 服务器安全通信的 PEM 编码证书和密钥文件。用于服务器与数据库之间的 TLS 连接。

关于生成自签名证书可参考官方说明：https://www.postgresql.org/docs/9.5/static/ssl-tcp.html，需要注意的是，自签名证书仅用于测试目的，不应在生产环境中使用。

有关在PostgreSQL服务器上配置SSL的更多详细信息，请参阅以下PostgreSQL文档：https://www.postgresql.org/docs/9.4/static/libpq-ssl.html
```
##### 配置 MySQL
如果使用 MySQL 数据库，则需要在 Fabric CA 服务器端的配置文件进行如下设置：
```
db:
  type: mysql
  datasource: root:rootpw@tcp(localhost:3306)/fabric_ca?parseTime=true&tls=custom
```
如果通过 TLS 连接到 MySQL 服务器，则还需要配置 db.tls.client 部分。如 PostgreSQL 的部分所述。

mySQL 数据库名称中允许使用字符限制。请参考：https://dev.mysql.com/doc/refman/5.7/en/identifiers.html

关于 MySQL 可用的不同模式，请参阅：https://dev.mysql.com/doc/refman/5.7/en/sql-mode.html，为正在使用的特定MySQL版本选择适当的设置。

#### 配置LDAP
LDAP（Lightweight Directory Access Protocol）：轻量目录访问协议。

Fabric CA服务器可以通过服务器端的配置连接到指定LDAP服务器。之后可以执行以下操作：

在注册之前读取信息进行验证
对用于授权的标识属性值进行验证
修改 Fabric CA 服务器的配置文件中的LDAP部分：
```
ldap:
   enabled: false
   url: <scheme>://<adminDN>:<adminPassword>@<host>:<port>/<base>
   userfilter: <filter>
   attribute:
      names: <LDAPAttrs>
      converters:
        - name: <fcaAttrName>
          value: <fcaExpr>
      maps:
        <mapName>:
            - name: <from>
              value: <to>
```
配置信息中各部分解释如下：
```
scheme：为 ldap 或 ldaps；
adminDN：是admin用户的唯一名称；
adminPassword：是admin用户的密码；
host：是LDAP服务器的主机名或IP地址；
port：是可选的端口号，默认 LDAP 为 389 ； LDAPS 为 636 ；
base：用于搜索的LDAP树的可选根路径；
filter：将登录用户名转换为可分辨名称时使用的过滤器；
LDAPAttrs：是一个LDAP属性名称数组，代表用户从LDAP服务器请求；

attribute.converters：部分用于将LDAP属性转换为结构CA属性，其中 fcaAttrName 是结构CA属性的名称; fcaExpr 是一个表达式。例如，假设是[“uid”]，是'hf.Revoker'，而是'attr（“uid”）=〜“revoker *”'。这意味着代表用户从LDAP服务器请求名为“uid”的属性。如果用户的'uid'LDAP属性的值以 revoker 开头，则为 hf.Revoker 属性赋予用户 true 的值；否则，为 hf.Revoker 属性赋予用户 false 的值。

attribute.maps：部分用于映射LDAP响应值。典型的用例是将与LDAP组关联的可分辨名称映射到标识类型。
```
配置好 LDAP 后，用户注册的过程如下：

Fabric CA 客户端或客户端 SDK 使用基本授权标头发送注册请求。
Fabric CA 服务器接收注册请求，解码授权头中的身份名称和密码，使用配置文件中的 “userfilter” 查找与身份名称关联的 DN（专有名称），然后尝试 LDAP 绑定用户身份的密码。如果 LDAP 绑定成功，则注册被通过。

## Fabric CA的具体使用
###Fabric CA 客户端命令
fabric-ca-client 命令可以与服务端进行交互, 包括五个子命令:
```
enroll：注册获取ECert
register：登记用户
getcainfo：获取CA服务的证书链
reenroll：重新注册
revoke：撤销签发的证书身份
version：Fabric CA 客户端版本信息
```
这些命令在执行时都是通过服务端的 RESTful 接口来进行操作的。

#### 注册用户
打开一个新的终端，首先，设置 fabric-ca-client 所在路径，然后设置 Fabric CA 客户端主目录。通过调用在 7054 端口运行的 Fabric CA 服务器来注册 ID 为 admin 且密码为 pass 的标识。
```
$ export PATH=$PATH:$GOPATH/bin
$ export FABRIC_CA_CLIENT_HOME=$HOME/fabric-ca/clients/admin
$ fabric-ca-client enroll -u http://admin:pass@localhost:7054
```
如果名称与密码不匹配， 则运行注册命令可能会产生如下错误：
```
Error: Response from server: Error Code: 20 - Authorization failure
```
解决方式: 
删除生成的目录，之后使用启动服务时的用户名与密码注册

或
返回至目录下重新启动服务， 然后在新终端中使用 admin:pass 注册
```
    $ cd ~
    $ fabric-ca-server start -b admin:pass
    打开新终端
    $ export PATH=$PATH:$GOPATH/bin
    $ export FABRIC_CA_CLIENT_HOME=$HOME/fabric-ca/clients/admin
    $ fabric-ca-client enroll -u http://admin:pass@localhost:7054
```
参数解释：
```
-u：进行连接的 fabric-ca-server 服务地址。
enroll 命令访问指定的 Fabric CA 服务，采用 admin 用户进行注册。 在 Fabric CA 客户端主目录下创建配置文件 fabric-ca-clien-config.yaml 和 msp 子目录，存储注册证书（ECert），相应的私钥和 CA 证书链 PEM 文件。我们可以在终端输出中看到指示 PEM 文件存储位置的相关信息。
```
生成的文件结构如下所示：
```
$ tree fabric-ca/clients/
fabric-ca/clients/
└── admin
    ├── fabric-ca-client-config.yaml
    └── msp
        ├── cacerts
        │   └── localhost-7054.pem
        ├── keystore
        │   └── 7441dddf832b4495cac12c05cc20b242f2ce545c5720010a83c11437157ac69d_sk
        ├── signcerts
        │   └── cert.pem
        └── user
提示：可以使用 $ tree fabric-ca/clients/ 命令查看目录结构
```
#### 登记用户
注册成功后的用户可以使用 register 命令来发起登记请求：

Fabric CA 服务器在注册期间进行了三次授权检查：

注册者（即调用者）必须具有 “hf.Registrar.Roles” 属性，其中包含逗号分隔的值列表，其中一个值等于要注册的身份类型； 如，如果注册商具有值为 “peer，app，user” 的 “hf.Registrar.Roles” 属性，则注册商可以注册 peer，app 和 user 类型的身份，但不能注册 orderer。
注册者的登记其范围内的用户。例如，具有 “a.b” 的从属关系的注册者可以登记具有 “a.b.c” 的从属关系的身份，但是可以不登记具有 “a.c” 的从属关系的身份。如果登记请求中未指定任何从属关系，则登记的身份将被授予注册者同样的归属范围。
如果满足以下所有条件，注册者可以指定登记用户属性：
注册者可以登记具有前缀 “hf” 的 Fabric CA 保留属性。只有当注册商拥有该属性并且它是hf.Registrar.Attributes 属性的值的一部分时。此外，如果属性是类型列表，则登记的属性值必须等于注册者具有的值的一个子集。如果属性的类型为 boolean，则只有当注册者的属性值为 “true” 时，注册者才能登记该属性。
注册自定义属性（即名称不以 'hf.' 开头的任何属性）要求注册者具有 'hf.Registar.Attributes' 属性，其中包含要注册的属性或模式的值。唯一支持的模式是末尾带有 “” 的字符串。例如，“a.b.\” 是匹配以 “a.b” 开头的所有属性名称的模式。例如，如果注册者具有hf.Registrar.Attributes = orgAdmin，则注册者可以在身份中添加或删除唯一的 orgAdmin 属性。
如果请求的属性名称为 “hf.Registrar.Attributes”，则执行附加检查以查看此属性的请求值是否等于 “hf.Registrar.Attributes” 的注册者值的子集。如，如果注册者的 hf.Registrar.Attributes 的值是 'a.b.，x.y.z' 并且所请求的属性值是 'a.b.c，x.y.z'，那么它是有效的，因为 'a.b.c' 匹配 'a.b '，'x.y.z' 匹配注册者的 'x.y.z' 值。
如下命令，使用管理员标识的凭据注登记 ID 为 “admin2” 的新用户，从属关系为 “org1.department1”，名为 “hf.Revoker” 的属性值为 “true”，以及属性名为 “admin”的值为 “true”。“：ecert” 后缀表示默认情况下，“admin” 属性及其值将插入用户的注册证书中，实现访问控制决策。
```
$ export FABRIC_CA_CLIENT_HOME=$HOME/fabric-ca/clients/admin
$ fabric-ca-client register --id.name admin2 --id.affiliation org1.department1 --id.attrs 'hf.Revoker=true,admin=true:ecert'
```
执行后输出：
```
Configuration file location: /home/kevin/.fabric-ca-client/fabric-ca-client-config.yaml
Password: KwnOlOhpfVit
```
命令执行成功后返回该新登记用户的密码。

如果想使用指定的密码, 在命令中添加选项 --id.secret password 即可

登记时可以将多个属性指定为 -id.attrs 标志的一部分，每个属性必须以逗号分隔。对于包含逗号的属性值，必须将该属性封装在双引号中。如：
```
$ fabric-ca-client register -d --id.name admin2 --id.affiliation org1.department1 --id.attrs '"hf.Registrar.Roles=peer,user",hf.Revoker=true'
```
#### 登记注册节点
登记Peer或Orderer节点的操作与登记用户身份类似；可以通过 -M 指定本地 MSP 的根路径来在其下存放证书文件

下面我们登记一个名为 peer1 的节点，登记时指定密码，而不是让服务器为生成。

登记节点:
```
$ export FABRIC_CA_CLIENT_HOME=$HOME/fabric-ca/clients/admin
$ fabric-ca-client register --id.name peer1 --id.type peer --id.affiliation org1.department1 --id.secret peer1pw
注册节点
$ export FABRIC_CA_CLIENT_HOME=$HOME/fabric-ca/clients/peer1
$ fabric-ca-client enroll -u http://peer1:peer1pw@localhost:7054 -M $FABRIC_CA_CLIENT_HOME/msp
```
参数说明：

-M： 指定生成证书存放目录 MSP 的路径, 默认为 "msp"
命令执行成功后会在 $FABRIC_CA_CLIENT_HOME 目录下生成指定的 msp 目录, 在此目录下生成 msp 的私钥和证书。

#### 其它命令
getcainfo
通常，MSP 目录的 cacerts 目录必须包含其他证书颁发机构的证书颁发机构链，代表 Peer 的所有信任根。

以下内容将在 localhost上启动第二个 Fabric CA 服务器，侦听端口 7055，名称为 “CA2”。这代表完全独立的信任根，并由区块链上的其他成员管理
```
$ export PATH=$PATH:$GOPATH/bin
$ export FABRIC_CA_SERVER_HOME=$HOME/ca2
$ fabric-ca-server start -b admin:ca2pw -p 7055 -n CA2
```
打开一个新终端，使用如下命令将CA2的证书链安装到peer1的MSP目录中
```
$ export PATH=$PATH:$GOPATH/bin
$ export FABRIC_CA_CLIENT_HOME=$HOME/fabric-ca/clients/peer1
$ fabric-ca-client getcainfo -u http://localhost:7055 -M $FABRIC_CA_CLIENT_HOME/msp
```
reenroll命令
如果注册证书即将过期或已被盗用。可以使用 reenroll 命令以重新生成新的签名证书材料
```
$ export FABRIC_CA_CLIENT_HOME=$HOME/fabric-ca/clients/peer1
$ fabric-ca-client reenroll
```
revoke命令
身份或证书都可以被撤销，撤销身份会撤销其所拥有的所有证书，并且还将阻止其获取新证书。被撤销后，Fabtric CA 服务器从此身份收到的所有请求都将被拒绝。

使用 revoke 命令的客户端身份必须拥有足够的权限（hf.Revoker为true, 并且被撤销者机构不能超出撤销者机构的范围）
```
$ export FABRIC_CA_CLIENT_HOME=$HOME/fabric-ca/clients/admin
$ fabric-ca-client revoke -e peer1 -r "affiliationchange"
```
参数说明：

-e：指定被撤销的身份
-r：指定被撤销的原因
命令执行后输出内容如下：
```
Configuration file location: /home/kevin/fabric-ca/clients/admin/fabric-ca-client-config.yaml
Sucessfully revoked certificates: [{Serial:21ed80434dd59cb1f80f89b85ebf55b3f677a54e AKI:1a99482cc8fe46349f0bd7ad7095985177708207} {Serial:4cf57dc2a8a70609e6eaaf3094e1ab3ff6aabe91 AKI:1a99482cc8fe46349f0bd7ad7095985177708207}]
```
另一种撤销身份的方式是可以指定其AKI（授权密钥标识符）和序列号来操作：
```
fabric-ca-client revoke -a xxx -s yyy -r <reason>
```
可以使用 openssl 命令获取 AKI 和证书的序列号，并将它们传递给 revoke 命令以撤销所述证书，如下所示：
```
serial=$(openssl x509 -in userecert.pem -serial -noout | cut -d "=" -f 2)
aki=$(openssl x509 -in userecert.pem -text | awk '/keyid/ {gsub(/ *keyid:|:/,"",$1);print tolower($0)}')
fabric-ca-client revoke -s $serial -a $aki -r affiliationchange
```
### 查看AKI和序列号
AKI: 公钥标识号, 代表了对该证书进行签发机构的身份

查看根证书的AKI与序列号信息:
```
$ openssl x509 -in $FABRIC_CA_CLIENT_HOME/msp/signcerts/cert.pem -text -noout
```
输出内容如下:
```
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:    # 序列号
            74:48:88:33:70:1a:01:a0:ad:32:29:6e:c5:ab:5a:fa:3b:91:25:a4
   ......
        X509v3 extensions:
           ......
            X509v3 Authority Key Identifier:     # keyid后面的内容就是 AKI
                keyid:45:B1:50:B6:CD:8A:8D:C5:9B:9E:5F:75:15:47:D6:C0:AD:75:FE:71

    ......
```
#### 单独获取AKI
```
$ openssl x509 -in $FABRIC_CA_CLIENT_HOME/msp/signcerts/cert.pem -text -noout | awk '/keyid/ {gsub (/ *keyid:|:/,"",$1);print tolower($0)}'
```
输出内容如下:
```
1a99482cc8fe46349f0bd7ad7095985177708207
```
#### 单独获取序列号
```
$ openssl x509 -in $FABRIC_CA_CLIENT_HOME/msp/signcerts/cert.pem -serial 
-noout | cut -d "=" -f 2
```
输出内容如下:
```
4CF57DC2A8A70609E6EAAF3094E1AB3FF6AABE91
```
