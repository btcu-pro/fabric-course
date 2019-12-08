package main

import (
	"fmt"
	web "github.com/kongyixueyuan.com/kongyixueyuan/hfsdkgoweb"
	"github.com/kongyixueyuan.com/kongyixueyuan/hfsdkgoweb/controller"
	"github.com/kongyixueyuan.com/kongyixueyuan/sdkInit"
	"github.com/kongyixueyuan.com/kongyixueyuan/service"
	"os"
)

const (
	configFile  = "config.yml"
	initialized = false
	SimpleCC    = "simplecc"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID:     "kevinkongyixueyuan",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/kongyixueyuan.com/kongyixueyuan/fixtures/artifacts/channel.tx",

		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",

		ChaincodeID:     SimpleCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/kongyixueyuan.com/kongyixueyuan/chaincode/",
		UserName:        "User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)

	//===========================================//

	serviceSetup := service.ServiceSetup{
		ChaincodeID: SimpleCC,
		Client:      channelClient,
	}

	msg, err := serviceSetup.SetInfo("hanxiaodong", "kongyixueyuan")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}

	//===========================================//

	app := controller.Application{
		Fabric: &serviceSetup,
	}
	web.WebStart(&app)
}
