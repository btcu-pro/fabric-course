package main

import (
"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/protos/peer"
"fmt"
"strconv"
)

type PaymentChaincode struct {

}

func main(){
    err := shim.Start(new(PaymentChaincode))
    if err != nil{
        fmt.Printf("启动 PaymentChaincode 时发生错误: %s", err)
    }
}

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


// peer chaincode query -n pay -C myc -c '{"Args":["find", "a"]}'
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