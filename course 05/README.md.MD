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
[github.com/chaincode/chaincode_example02/go/](./files/chaincode_example02.go)

