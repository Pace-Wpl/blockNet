package daemontask

import (
	"fmt"
	"net"
	"time"

	"github.com/blockNet/clientTcp/def"
	"github.com/blockNet/clientTcp/service"
	ch "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type Task struct {
	Controller chan string
	Conn       *net.TCPConn
}

func (t *Task) eventResult(notifier <-chan *fab.CCEvent) bool {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
		return true
	case e := <-t.Controller:
		fmt.Println("收到结束码:%s,正在结束...\n", e)
		return false
	}
}

func regitserEvent(client *ch.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Printf("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func (t *Task) ListenTask(event string, CarNum string) {
	for {
		reg, notifier := regitserEvent(service.ServiceClient.Client, service.ServiceClient.ChaincodeID, event)

		if t.eventResult(notifier) {
			switch event {
			case def.LOCK_EVENT:
				service.LockCheck(CarNum, t.Conn)
			case def.FAULTCODE_EVENT:
				service.FaulCheck(CarNum, t.Conn)
			case def.RGL_EVENT_SPPED:
				service.RglCheck(CarNum, t.Conn)
			}
		} else {
			service.ServiceClient.Client.UnregisterChaincodeEvent(reg)
			return
		}

		service.ServiceClient.Client.UnregisterChaincodeEvent(reg)
	}
}

//开启监听
func (t *Task) StartDaemon(CarNum string) {
	go t.ListenTask(def.FAULTCODE_EVENT, CarNum)
	go t.ListenTask(def.LOCK_EVENT, CarNum)
	go t.ListenTask(def.RGL_EVENT_SPPED, CarNum)
}

//停止监听
func (t *Task) StopDaemon() {
	t.Controller <- "done1"
	time.Sleep(time.Duration(2) * time.Second)
	t.Controller <- "done2"
	time.Sleep(time.Duration(2) * time.Second)
	t.Controller <- "done3"
	time.Sleep(time.Duration(2) * time.Second)
}
