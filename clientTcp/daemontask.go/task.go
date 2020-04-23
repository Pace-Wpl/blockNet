package daemontask

import (
	"fmt"
	"net"

	"github.com/blockNet/clientTcp/def"
	"github.com/blockNet/clientTcp/service"
	ch "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type Task struct {
	Controller chan string
}

func eventResult(notifier <-chan *fab.CCEvent) bool {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
		return true
	}
}

func regitserEvent(client *ch.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Printf("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func (t *Task) ListenTask(conn *net.TCPConn, event string, carNum string) {
	for {
		reg, notifier := regitserEvent(service.ServiceClient.Client, service.ServiceClient.ChaincodeID, event)

		if eventResult(notifier) {
			switch event {
			case def.LOCK_EVENT:
				service.LockCheck(carNum)
			case def.FAULTCODE_EVENT:
				service.FaulCheck(carNum)
			case def.RGL_EVENT_SPPED:
				service.RglCheck(carNum)
			}
		}

		service.ServiceClient.Client.UnregisterChaincodeEvent(reg)
	}
}

func StartDaemon(conn *net.TCPConn, carNum string) {
	t := &Task{}
	go t.ListenTask(conn, def.FAULTCODE_EVENT, carNum)
	go t.ListenTask(conn, def.LOCK_EVENT, carNum)
	go t.ListenTask(conn, def.RGL_EVENT_SPPED, carNum)
}
