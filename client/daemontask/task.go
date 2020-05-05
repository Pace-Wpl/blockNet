package daemontask

import (
	"fmt"
	"time"

	"github.com/blockNet/client/def"
	"github.com/blockNet/client/service"
	ch "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type Task struct {
	Controller chan string
	// Conn       *net.TCPConn
	Flag bool
}

func (t *Task) eventResult(notifier <-chan *fab.CCEvent) bool {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
		return true
	case e := <-t.Controller:
		fmt.Println("收到结束码:%s,监听正在结束...\n", e)
		return false
	}
}

func regitserEvent(client *ch.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册监听失败,请重新开启...")
	}
	fmt.Println("注册监听：" + eventID)
	// fmt.Println("正在监听...")
	return reg, notifier
}

func (t *Task) ListenTask(event string, CarNum string) {
	for {
		reg, notifier := regitserEvent(service.ServiceClient.Client, service.ServiceClient.ChaincodeID, event)

		if t.eventResult(notifier) {
			switch event {
			case CarNum + def.LOCK_EVENT:
				service.LockCheck(CarNum)
			// case def.FAULTCODE_EVENT:
			// 	service.FaulCheck(CarNum, t.Conn)
			case CarNum + def.RGL_EVENT_SPPED:
				service.RglCheck(CarNum)
			case CarNum + def.COLLIS_EVENT:
				service.CollisCheck(CarNum)
				// case def.ON_ROAD_EVENT:
				// 	service.OnRoadListen(CarNum, t.Conn)
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
	t.Flag = true
	go t.ListenTask(CarNum+def.COLLIS_EVENT, CarNum)
	go t.ListenTask(CarNum+def.LOCK_EVENT, CarNum)
	go t.ListenTask(CarNum+def.RGL_EVENT_SPPED, CarNum)
	// go t.ListenTask(def.ON_ROAD_EVENT, CarNum)
}

//停止监听
func (t *Task) StopDaemon() {
	if t.Flag {
		t.Controller <- "done1"
		time.Sleep(time.Duration(2) * time.Second)
		t.Controller <- "done2"
		time.Sleep(time.Duration(2) * time.Second)
		t.Controller <- "done3"
		time.Sleep(time.Duration(2) * time.Second)
		// t.Controller <- "done4"
		// time.Sleep(time.Duration(2) * time.Second)
	}
}
