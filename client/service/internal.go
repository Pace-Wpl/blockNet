package service

import (
	"bytes"
	"fmt"
	"time"

	ch "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

//用eventID注册链码事件
func regitserEvent(client *ch.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	fmt.Println("足册:" + eventID)
	if err != nil {
		fmt.Printf("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

//用eventID接受链码事件结果
func eventResult(notifier <-chan *fab.CCEvent, eventID string) (string, error) {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
		if bytes.Equal(ccEvent.Payload, []byte{}) {
			return "Successful!", nil
		} else {
			return string(ccEvent.Payload) + "11", nil
		}
	case <-time.After(time.Second * 20):
		fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
		return "不能根据指定的事件ID接收到相应的链码事件(" + eventID + ")\n", nil
	}
}
