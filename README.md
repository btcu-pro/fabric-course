# BTCU_Fabric

## Constents
- [BTCU_Fabric](#btcufabric)
  - [Constents](#constents)
  - [0. Fabric 简介](#0-fabric-%e7%ae%80%e4%bb%8b)
    - [超级账本背景](#%e8%b6%85%e7%ba%a7%e8%b4%a6%e6%9c%ac%e8%83%8c%e6%99%af)
    - [Fabric 背景](#fabric-%e8%83%8c%e6%99%af)
    - [Fabric 架构简介](#fabric-%e6%9e%b6%e6%9e%84%e7%ae%80%e4%bb%8b)
  - [1. 开发环境搭建（Docker、Mac、Windows、Linux）](#1-%e5%bc%80%e5%8f%91%e7%8e%af%e5%a2%83%e6%90%ad%e5%bb%badockermacwindowslinux)
    - [Docker 简介](#docker-%e7%ae%80%e4%bb%8b)
    - [通过 Windows 下载 Docker 镜像来得到 Fabric 环境](#%e9%80%9a%e8%bf%87-windows-%e4%b8%8b%e8%bd%bd-docker-%e9%95%9c%e5%83%8f%e6%9d%a5%e5%be%97%e5%88%b0-fabric-%e7%8e%af%e5%a2%83)
      - [安装 Docker 环境](#%e5%ae%89%e8%a3%85-docker-%e7%8e%af%e5%a2%83)
      - [在 Docker 里下载定制好了的 Fabric 镜像：](#%e5%9c%a8-docker-%e9%87%8c%e4%b8%8b%e8%bd%bd%e5%ae%9a%e5%88%b6%e5%a5%bd%e4%ba%86%e7%9a%84-fabric-%e9%95%9c%e5%83%8f)
    - [Linux 环境下直接安装](#linux-%e7%8e%af%e5%a2%83%e4%b8%8b%e7%9b%b4%e6%8e%a5%e5%ae%89%e8%a3%85)
      - [安装 Golang](#%e5%ae%89%e8%a3%85-golang)

## 0. Fabric 简介

### 超级账本背景

超级账本(Hyperledger)项目是首个面向企业应用场景的开源分布式账本平台。

2015 年 12 月，由开源世界的旗舰组织 Linux 基金会牵头，30 家初始企业成员（包括 IBM、Intel、摩根大通、思科、R3 等），共同宣布了 Hyperledger 联合项目成立。超级账本项目为透明、公开、去中心化的企业级分布式账本技术提供开源参考实现，并推动区块链和分布式账本相关协议、规范和标准的发展。

作为一个联合项目（collaborative project），超级账本由面向不同目的和场景的子项目构成。目前包括 Fabric、Sawtooth、Iroha、Blockchain Explorer、Cello、Indy、Composer、Burrow 等 8 大顶级项目。

这次课程将专注于 Fabric。

### Fabric 背景

Fabric 是最早加入到超级账本项目中的顶级项目，Fabric 由 IBM、DAH 等企业于 2015 年底提交到社区。该项目的定位是面向企业的分布式账本平台，创新地引入了权限管理支持，设计上支持可插拔、可扩展，是首个面向联盟链场景的开源项目。
Fabric 基于 Go 语言实现，目前已经发布了 1.4 版本，同时包括 Fabric CA、Fabric SDK 等多个子项目。

GitHub 地址: https://github.com/hyperledger/fabric/

### Fabric 架构简介

4 种不同种类的服务节点：
* 背书节点（Endorser）：负责对交易的提案（proposal）进行检查和背书，计算交易执行结果；
* 确认节点（Committer）：负责在接受交易结果前再次检查合法性，接受合法交易对账本的修改，并写入区块链结构；
* 排序节点（Order）：对所有发往网络中的交易进行排序，将排序后的交易按照配置中的约定整理为区块，之后提交给确认节点进行处理；
* 证书节点（CA）：负责对网络中所有的证书进行管理，提供标准的 PKI 服务。

一笔交易的典型流程图：
![fabric 交易流程](image/fabric02.jpg)

2 种通道：
* 系统通道（system channel）唯一，独立。负责管理网络中的各种配置信息，并完成对其他应用通道（application channel）的创建。
* 应用通道（application channel）可以有多个。供用户发送交易使用。

启动一个 Fabric 网络的主要步骤：
1. 预备网络内各项配置，包括网络中成员的组织结构和对应的身份证书（使用 cryptogen 工具完成）；生成系统通道的初始配置区块文件，新建应用通道的配置更新交易文件以及可能需要的锚节点配置更新交易文件（使用 configtxgen 工具完成）。
2. 使用系统通道的初始配置区块文件启动排序节点，排序节点启动后自动按照指定配置创建系统通道。
3. 不同的组织按照预置角色分别启动 Peer 节点。这个时候网络不存在应用通道，Peer 节点也并没有加入网络中。
4. 使用新建应用通道的配置更新交易文件，向系统通道发送交易，创建新的应用通道。
5. 让对应的 Peer 节点加入所创建的应用通道中，此时 Peer 节点加入网络，可以转变接受交易了。
6. 用户通过客户端向网络中安装注册链码（chaincode），链码容器启动成功后用户即可对链码进行调用，将交易发送到网络中去。

## 1. 开发环境搭建（Docker、Mac、Windows、Linux）

### Docker 简介

[Docker](https://www.docker.com/) 是一类虚拟化技术，类似虚拟机，安装好 Docker 这个软件后，可以在 Docker 里面运行一些定制好了的镜像，比如今天我们用到的就是在 Linux 操作系统里安装了 Fabric 环境的镜像，这样就省去了直接安装 Fabric 的麻烦。更多的镜像可以参考 [Docker Hub](https://hub.docker.com/)。然后本身 Docker 这个软件支持 Mac、Windows 和 Linux，因此也是一个很好的跨平台工具。

### 通过 Windows 下载 Docker 镜像来得到 Fabric 环境

参考：https://docs.docker.com/docker-for-windows/install/ 
（参考表示下面的具体安装文档是参考这个链接 + 实际情况得到的，因此一般来讲可以直接根据下面的步骤操作即可，下同）

#### 安装 Docker 环境

从官网这个[链接](https://hub.docker.com/search?q=&type=edition&offering=community)进去，然后下拉，找到 Docker Desktop for Windows，点击下载。
![下载Docker软件](./image/fabric03.png)

这里显示需要先登录：
![需要登录](./image/fabric04.png)

如果之前有账户的直接登录，没有的就点击 Sign Up：
![登录界面](image/fabric05.png)

登录之后就可以直接下载了（建议把下载链接放在迅雷里面，在浏览器里面下载 10k 左右，迅雷里面 10M）：
![下载界面](image/fabric06.png)

安装时发现需要 Windows 10 的专业板或者企业版本：
![家庭版无法安装](image/fabric08.png)

然后查询发现本系统是家庭版：
![Windows版本](image/fabric07.png)

因此需要安装 docker-toolbox:
在 https://github.com/docker/toolbox/releases ：
![下载 toolbox](image/fabric09.png)

然后安装 docker toolbox 就是按照默认一路确认下去。过程中可能需要安装几个其他软件，同样都是确认。最后得到成功安装了的 Docker Quickstart:  
![Docker Quickstart](image/fabric10.png)

双击后启动 docker，这过程需要下载一个镜像，需要一些时间最后成功安装，运行 `docker --version` ，显示版本即成功安装 docker ：
![start docker](image/fabric11.png)

总结：Windows 10 安装 Docker 有两种情况，如果版本是 Windows 10 专业版或者企业版，可以直接通过 `Docker for Windows Installer.exe` 安装，否则可以通过 `DockerToolbox-19.03.1.exe` 安装。

#### 在 Docker 里下载定制好了的 Fabric 镜像：

Fabric 有多个镜像，下面是对应的依赖关系：  
![images](image/fabric12.png)  

这里需要安装的是 fabric-peer, fabric-orderer, fabric-ca, fabric-tools, fabric-ccenv。  
更多的镜像参考：  
https://hub.docker.com/search?q=hyperledger&type=image 
![image_web](image/fabric13.png)  
以下命令直接下载对应 fabric 最新的镜像，也就是 fabric 1.4 版本。
```shell
docker pull hyperledger/fabric-peer \
    && docker pull hyperledger/fabric-orderer \
    && docker pull hyperledger/fabric-ca \
    && docker pull hyperledger/fabric-tools \
    && docker pull hyperledger/fabric-ccenv
```


下载完成后，用 `docker images` 查看，可以看到刚刚下载的 5 个镜像。
![pull_linux](image/fabric15.png)   
到这里，Windows 下的 Docker 环境配置好了

### Linux 环境下直接安装

#### 安装 Golang
从官网下载最新版本：
```shell
curl -O https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz
```
解压：
```shell
tar -xvf go1.12.9.linux-amd64.tar.gz
```
得到 `./go/` 文件夹：
```shell
$ ls go/
api      bin              CONTRIBUTORS  favicon.ico  LICENSE  PATENTS  README.md   src   VERSION
AUTHORS  CONTRIBUTING.md  doc           lib          misc     pkg      robots.txt  test
```
这个文件夹里面就有 go 语言的配套环境了，然后设置当前用户的环境变量。

用编辑器打开 ~/.bashrc 文件，比如我是用 Emacs：
```shell
$ emacs ~/.bashrc
```
在最后一行添加：
```shell
export PATH=$PATH:/home/flyq/workspaces/golang/go/bin/
```
主要，添加的这行每个人的路径不同，因此这行代码也不同，如下图，需要根据自己电脑环境对应目录的路径得到：
![dir](image/fabric16.png)
![dir2](image/fabric17.png)

然后保存好，更新一下：
```shell 
source ~/.bashrc
```

运行`go version`出现以下结果即表示安装成功：
```shell
$ go version 
go version go1.12.9 linux/amd64
```

最后设置一下 GOPATH 环境变量，同样是修改 `~/.bashrc` 文件：
创建一个新建目录，并指定它是 GOPATH：
![gopath](image/fabric18.png)


然后保存好，更新一下：
```shell 
source ~/.bashrc
```
go 环境已经安装并配置好了。


