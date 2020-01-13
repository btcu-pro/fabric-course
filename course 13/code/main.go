package main

import (
	"fmt"
	web "github.com/btcu-pro/btcu/hfsdkgoweb"
	"github.com/btcu-pro/btcu/hfsdkgoweb/controller"
	"github.com/btcu-pro/btcu/sdkInit"
	"github.com/btcu-pro/btcu/service"
	"os"
)

const (
	configFile  = "config.yml"
	initialized = false
	SimpleCC    = "simplecc"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID:     "demobtcu",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/btcu-pro/btcu/fixtures/artifacts/channel.tx",

		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer.demo.btcu.com",

		ChaincodeID:     SimpleCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/btcu-pro/btcu/chaincode/",
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

	msg, err := serviceSetup.SetInfo("hanxiaodong", "btcu")
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
