package service

import (
	"fmt"
	"testing"

	"github.com/blockNet/clientTcp/def"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestSdkInitWorkFlow(t *testing.T) {
	t.Run("sdk setup all \n", testSdkSetUp)
}

func testSdkSetUp(t *testing.T) {
	sdk, err := SetupSDK("/home/pace/go/src/github.com/blockNet/client/config.yaml", false)
	if err != nil {
		panic(err.Error())
	}
	defer sdk.Close()

	info := &def.InitInfo{ChannelID: def.CHAINNEL_ID, UserName: def.USER_NAME, OrgName: def.ORG_NAME}

	channelClient, err := GetChannelClient(sdk, info)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(channelClient)

}
