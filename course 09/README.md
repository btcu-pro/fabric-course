# Course 09

gossip协议是一个神奇的协议。它常用于P2P的通信协议，这个协议就是模拟人类中传播谣言的行为而来。简单的描述下这个协议，首先要传播谣言就要有种子节点。种子节点每秒都会随机向其他节点发送自己所拥有的节点列表，以及需要传播的消息。任何新加入的节点，就在这种传播方式下很快地被全网所知道。这个协议的神奇就在于它从设计开始就没想到信息一定要传递给所有的节点，但是随着时间的增长，在最终的某一时刻，全网会得到相同的信息。

Hyperledger Fabric通过拆分工作量为交易执行节点（背书(Endorsing)和提交(Committing)）和交易排序节点，来优化区块链网络的性能，安全性和可伸缩性。这种网络操作的解耦需要一个安全的、可靠的和可伸缩的数据传播协议，以确保数据的完整性和一致性。为了满足这些要求，Hyperledger Fabric实现 gossip数据传播协议 。
## Gossip 协议
### Gossip 协议的概念
节点(Peer)利用gossip以可扩展的方式广播账本(Ledger)和频道(Channel)的数据。Gossip的消息发送是连续的，频道(Channel)上的每个节点(Peer)都不断地从多个其他节点(Peer)接收最新的和一致的账本数据。每个gossip消息都会被签名，从而，拜占庭参与者(Byzantine Participants)发送的伪造消息将会很容易地被识别出来，并且，阻止这些消息被分发到不想要的目标。节点(Peer)会受到网络延时、网络分区或其他因素的影响，从而导致丢失区块(Block)，但最终会通过联系那些拥有缺失区块的节点(Peer)，而同步到最新的帐本状态。

基于gossip的数据传播协议在 Hyperledger Fabric 网络中，主要执行三项功能：

1. 通过持续不断地识别可用的会员节点，并最终检测已离线的节点，来管理节点发现(peer discovery)和频道(在线)会员(channel membership)。
2. 在频道(Channel)中的所有节点(Peer)之间传播账本数据。可以识别与频道(Channel)中其余节点不同步的节点(Peer)，标示出丢失的数据块，并通过复制正确的数据来同步自身。
3. 通过以P2P的方式更新账本数据，使新连接的节点(Peer)加速。

Gossip的广播，首先，是由节点(Peer)接收频道(Channel)上的其他节点(Peer)消息，然后，将这些消息转发给频道(Channel)中的多个随机节点(Peer)，这个数量是一个可配置的常量。节点(Peer)也可以使用拉取(pull)机制，而不是等待信息的传递。这个循环不断地重复，最终，频道(在线)会员(channel membership)，账本(Ledger)和状态信息，将会持续不断地被同步至最新的结果。为了传播新的区块，该频道(Channel)中的 领导 节点(leader peer)将从排序服务(Ordering Service)中拉取(pull)数据，并开始向其他节点(Peer)传播gossip消息。

### Gossip的信息发送
在线节点(Peer)通过持续不断地广播“活着”的消息，来表明他们的可用性，每个消息都包含了 公钥基础设施(Public Key Infrastructure - PKI) ID和发件人针对消息的签名。节点(Peer)通过收集这些“活着”的消息来维护频道(在线)会员(channel membership)。如果没有节点(Peer)收到来自特定节点(Peer)的“活着”消息，则这个“死了”的节点将最终从频道(在线)会员(channel membership)中清除。因为“活着”消息是被加密签名的，并且，由于缺乏根证书授权机构（CA）的签名密钥，所以，恶意节点(Peer)永远不能冒充其他节点(Peer)。

除了自动转发收到的消息之外，频道(Channel)中的节点(Peer)之间，还会有一个状态核对过程，来同步 世界状态(world state) 。每个节点(Peer)持续不断地从该频道(Channel)的其他节点(Peer)中拉取(pull)区块，以便在发现差异的情况下修复自己的状态(state)。基于gossip的数据传播不需要保持固定的网络连接，因此，此过程非常可靠地保证了共享帐本的数据一致性和完整性，包括对节点崩溃的容错性。

由于频道(Channel)是相互隔离的，所以一个频道(Channel)上的节点(Peer)不能发送消息或分享信息给其他任何的频道(Channel)。虽然任意节点(Peer)可以从属于多个频道(Channel)，但通过应用基于节点频道订阅(peers’ channel subscription)的消息路由策略(message routing policy)，分区的消息发送机制(partitioned messaging)可以防止区块被传播到不在此频道(Channel)中的节点(Peer)。


### Gossip 的特点

1）扩展性

网络可以允许节点的任意增加和减少，新增加的节点的状态最终会与其他节点一致。

2）容错

网络中任何节点的宕机和重启都不会影响 Gossip 消息的传播，Gossip 协议具有天然的分布式系统容错特性。

3）去中心化

Gossip 协议不要求任何中心节点，所有节点都可以是对等的，任何一个节点无需知道整个网络状况，只要网络是连通的，任意一个节点就可以把消息散播到全网。

4）一致性收敛

Gossip 协议中的消息会以一传十、十传百一样的指数级速度在网络中快速传播，因此系统状态的不一致可以在很快的时间内收敛到一致。消息传播速度达到了 logN。

5）简单

Gossip 协议的过程极其简单，实现起来几乎没有太多复杂性。


### Gossip 类型

Gossip 有两种类型：

Anti-Entropy（反熵）：以固定的概率传播所有的数据
Rumor-Mongering（谣言传播）：仅传播新到达的数据
Anti-Entropy 是 SI model，节点只有两种状态，Suspective 和 Infective，叫做 simple epidemics。

Rumor-Mongering 是 SIR model，节点有三种状态，Suspective，Infective 和 Removed，叫做 complex epidemics。

其实，Anti-entropy 反熵是一个很奇怪的名词，之所以定义成这样，Jelasity 进行了解释，因为 entropy 是指混乱程度（disorder），而在这种模式下可以消除不同节点中数据的 disorder，因此 Anti-entropy 就是 anti-disorder。换句话说，它可以提高系统中节点之间的 similarity。

在 SI model 下，一个节点会把所有的数据都跟其他节点共享，以便消除节点之间数据的任何不一致，它可以保证最终、完全的一致。

由于在 SI model 下消息会不断反复的交换，因此消息数量是非常庞大的，无限制的（unbounded），这对一个系统来说是一个巨大的开销。



但是在 Rumor Mongering（SIR Model） 模型下，消息可以发送得更频繁，因为消息只包含最新 update，体积更小。而且，一个 Rumor 消息在某个时间点之后会被标记为 removed，并且不再被传播，因此，SIR model 下，系统有一定的概率会不一致。

而由于，SIR Model 下某个时间点之后消息不再传播，因此消息是有限的，系统开销小。



### Gossip 中的通信模式

在 Gossip 协议下，网络中两个节点之间有三种通信方式:

1. Push: 节点 A 将数据 (key,value,version) 及对应的版本号推送给 B 节点，B 节点更新 A 中比自己新的数据
2. Pull：A 仅将数据 key, version 推送给 B，B 将本地比 A 新的数据（Key, value, version）推送给 A，A 更新本地
3. Push/Pull：与 Pull 类似，只是多了一步，A 再将本地比 B 新的数据推送给 B，B 则更新本地

如果把两个节点数据同步一次定义为一个周期，则在一个周期内，Push 需通信 1 次，Pull 需 2 次，Push/Pull 则需 3 次。虽然消息数增加了，但从效果上来讲，Push/Pull 最好，理论上一个周期内可以使两个节点完全一致。直观上，Push/Pull 的收敛速度也是最快的。

## Fabric中数据同步的实现
### Hyperledger Fabric中的Gossip
Fabric 中 的各个 Peer 节点之间利用 Gossip 协议来完成区块广播以及状态同步的过程。Gossip 消息是连续的，通道上的每个 Peer 节点都不断地接收来自多个节点已完成一致性的区块数据。每条传输的 Gossip 消息都有相应的签名，从而由拜占庭参与者发送的伪造消息很容易地识别来，并且可以防止将消息分发给不在同一通道中的其它节点。受到延迟、网络分区或其他导致区块丢失的原因影响的节点，最终将通过联系已经拥有这些缺失区块的节点，与当前账本状态进行同步。

在 Hyperledger Fabric 网络中基于 Gossip 的数据传播协议在 Fabric 网络上执行三个主要功能：

1. 通过不断识别可用的成员节点并最终监测节点离线状态的方式，对节点的发现和通道中的成员进行管理。
2. 将分类帐本数据传播到通道上的所有节点。任何节点中如有缺失区块都可以通过从通道中其它节点复制正确的数据来标识缺失的区块并同步自身。
3. 在通道上的所有节点上同步分类帐状态。通过允许点对点状态传输更新账本数据，保证新连接的节点以最快的速度实现数据同步。

基于 gossip 的广播由节点接收来自通道内其他节点的消息，然后将这些消息转发给随机选择的且在同一通道内的若干个邻居节点，这种循环不断重复，使通道中所有的成员节点的账本和状态信息不断保持与当前的最新状态同步。对于新区块的传播，通道上的 Leader Peer 节点从 Ordering 服务中提取数据，并向随机选择的邻居节点发起 Gossip 传播。

随机选择的邻居节点数量可以通过配置文件进行配置声明。节点也可以使用拉取机制，而不是等待消息的传递。

客户端应用程序将交易提案请求提交给背书节点（Endorse Peer），背书节点处理并背书签名后返回响应，然后提交给 Ordering 服务进行排序，排序服务达成共识后生成区块，通过 deliver（）广播给各个组织中通过选举方式选择的作为代表能够连接到排序服务的 Leader Peer 节点，Leader Peer 节点随机选择N个节点将接收到的区块进行分发。另外，为了保持数据同步，每个节点会在后台周期性地与其它随机的N个节点的数据进行比较，以保持区块数据状态的同步。

### Leader节点选举
在 Hyperledger Fabric 网络中，每一个组织都会通过领导选举机制选择一个节点（Leader Peer），该节点将保持与 Ordering 服务的连接，并在其所在组织的节点之间分发从 Ordering 服务节点接收到的新区块。利用领导人选举为系统提供了有效利用 Ordering 服务带宽的能力。在 Hyperledger Fabric 中实现领导人选举有两种方式：

静态选举：由系统管理员手动配置实现，指定组织中的一个 peer 节点作为领导节点代表组织与 Ordering 服务建立连接。
动态选举：通过执行领导人选举程序，动态从组织中选择一个 peer 节点成为领导者节点，从Ordering 服务中拉出区块，并将块分发给组织中的其他 peer 节点。
静态选举
使用静态领导选举可以通过在配置文件中指定相关的参数来实现。可以定义一个节点为 Leader Peer，也可定义多个节点或组织内所有节点都为 Leader Peer。

实现静态选举机制，需要在 core.yaml 中配置以下参数:
```
peer:
    gossip:
        useLeaderElection: false    # 是否指定使用选举方式产式Leader
        orgLeader: true    # 是否指定当前节点为Leader
```
或者可以使用环境变量来配置和覆盖相应的参数：
```
export CORE_PEER_GOSSIP_USELEADERELECTION=false
export CORE_PEER_GOSSIP_ORGLEADER=true
```
注意，如果两个值全部都指定为 false， 那么代表 peer 节点不会成为领导者。

动态选举
动态领导选举可以在各个组织内各自动态选举一个 Leader 节点，它将代表各个连接到 Ordering 服务并拉出新的区块。

当选的 Leader 节点必须向组织内的其他节点定期发送心跳信息，作为处于活跃的证据。如果一名或多名节点在指定的一段时间内得不到最新消息，网络将启动新一轮领导人选举程序，最终选出新的 Leader 节点。

启用动态选举机制，需要在 core.yaml 中配置以下参数:
```
peer:
    gossip:
        useLeaderElection: true     # 是否指定使用选举方式产式Leader
        orgLeader: false    # 是否指定当前节点为Leader
```
或者，可以使用环境变量来配置和覆盖相应参数：
```
export CORE_PEER_GOSSIP_USELEADERELECTION=true
export CORE_PEER_GOSSIP_ORGLEADER=false
```
core.yaml 以下配置内容指定了动态选举 Leader 的相关信息：
```
peer:
    gossip:
         election:   # 选举Leader配置     
            startupGracePeriod: 15s       # 最长等待时间 
            membershipSampleInterval: 1s  # 检查稳定性的间隔时间     
            leaderAliveThreshold: 10s     # 进行选举的间隔时间
            leaderElectionDuration: 5s    # 声明自己为Leader的等待时间
```
锚节点（Anchor Peer）
锚节点主要用于启动来自不同组织的节点之间的 Gossip 通信。锚节点作为同一通道上的另一组织的节点的入口点，可以与目标锚节点所在组织中的每个节点通信。跨组织的 Gossip 通信必须包含在通道的范围内。

由于跨组织的通信依赖于 Gossip，某一个组织的节点需要知道来自其它组织的节点的至少一个地址(从这个节点，可以找到该组织中的所有节点的信息)。所以添加到通道的每个组织应将其节点中的至少一个节点标识为锚节点（也可以有多个锚节点，以防止单点故障）。网络启动后锚节点地址存储在通道的配置块中。

可以通过在 configtx.yaml 配置文件指定锚节点：
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

        AnchorPeers:    # 指定当前组织的锚节点
            - Host: peer0.org1.example.com
              Port: 7051
    ......
```

### Fabric的数据同步实现
Hyperledger Fabric 是一个分布式区块链网络，所有的 peer 节点都会保存共享分类帐的副本（即所有事务的确切历史记录）。当新区块产生后必须通过分布式网络，使分类帐的副本在所有节点之间保持同步。

在较高的层次上，该过程如下所示：

新的交易被提交给 Ordering 服务进行排序。
Ordering 服务在排序之后创建一个新区块（包含新的交易）。
Ordering 服务将新产生的区块交给所有 Peer。
但在Hyperledger Fabric 网络中实际发生的情况是，Ordering 服务只向每个组织中的单个节点（Leader Peer）提供新的区块。通过 Gossip 的过程， Peer 节点自己完成了将新区块传播到其它 Peer 节点的工作：

Peer 节点接收到新的消息。
该节点将消息发送到预先指定数量（随机选择的 Fabric中 默认为3个 Peer）的其他 Peer 节点。
接收到消息的每一个 Peer 节点再将消息转发给预定数量的其他 Peer 节点。
依此类推，直到所有的 Peer 节点都收到了新的消息。
上面的过程称之为广播，它是一种基于推送（Push-based）的方式，通过网络传输信息，Fabric 的 Gossip 系统使用它来向所有 Peer 节点分发消息。

Gossip 协议的关键组成部分是每个节点将消息随机选择并转发给网络中其它节点。这意味着每个节点都知道网络中的所有节点，因此可以在相应的 Peer 节点中进行选择。那么，某一个节点是如何知道组织内的所有节点呢？并且如果有 Peer 节点与网络断开连接并在后期重新连接，则它将错过广播过程。

在 Hyperledger Fabric 中，每个节点都会随机性的向预先定义数量的其它节点定期广播一条消息，指示它仍处于活动状态并连接到网络。每个节点都维护着自己的网络中所有节点的列表（处于活跃的节点和无响应的节点）。

当某一个节点 A 收到来自节点 B 的“活跃”消息时，它将节点 B 标记为“有效”（Peer B是网络中的一个有效节点）
如果过了一段时间，节点 A 没有收到来自节点 B 的“活跃”消息，Peer A 节点会定期尝试连接 Peer B 节点，确认是否真的无响应。如果无响应将节点 B 标记为“死亡”（Peer B不再是网络的有效节点）。
这种情况之下需要一个基于拉取（Pull-based）的实现机制来向其它 Peer 节点请求它丢失的数据。在Hyperledger Fabric中， Peer 节点之间定期相互交换成员资格数据（ Peer 节点列表，活动和死亡）和分类帐本数据（事务块）。在这种机制下， Peer 节点即使因为故障或其它原因导致错过了接收新区块的广播或因为其它原因产生了缺失区块，但仍然在加入网络之后可以与其它的 Peer 节点交换信息以保持数据同步。

### Fabric数据同步

正如上图所示，Hyperledger Fabric使用对等体之间的 Gossip 作为容错和可扩展机制，以保持区块链分类账的所有副本同步，它减少了Orderer 节点上的负载。由于不需要固定连接来维护基于Gossip的数据传播，因此该流程可以可靠地为共享账本保证数据的一致性和完整性，包括对节点崩溃的容错。

另外，某些节点可以加入多个不同的通道，但是通过将基于节点通道订阅的机制作为消息分发策略，由于通道之间实现了相互隔离，一个通道上的节点不能在其他通道上发送或共享信息，所以节点无法将区块传播给不在通道中的节点。

点对点消息的安全性由节点的TLS层处理，不需要签名。节点通过其由CA分配的证书进行身份验证。节点在Gossip层的身份认证会通过TLS证书体现。账本中的区块由排序服务进行签名，然后传递给通道中的领导者节点。

身份验证过程由节点的成员管理服务的提供者（MSP）进行管理。当节点第一次连接到通道中的时候，TLS会话将与成员身份绑定。这本质上是通过网络和通道中的成员身份对连接的每个节点进行身份验证。
