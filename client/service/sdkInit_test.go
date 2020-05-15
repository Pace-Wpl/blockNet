package service

import (
	"fmt"
	"testing"

	"github.com/blockNet/client/def"
)

func TestMain(m *testing.M) {
	m.Run()
}

// func BenchmarkRaft_TPS(b *testing.B) {

//     for i := 0; i < b.N; i++ {
//         b.ReportAllocs() // 这里可以直接调用 ReportAllocs 方法，就省去了再命令行中输入 -benchmem ，用于查看内存分配的大小和次数
//         _, _ = a.GetReport("devices", "appsinfo", "")
//     }
// }

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
