package heart

/*
	负责从kafka中读取消息，维护活跃节点
*/

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofish2020/easydispatch/heart/kafka"
)

// global variable
var addrList = make(map[string]time.Time)
var mutex sync.Mutex

// listen and revice heart beat
func ListenHeartbeat() {
	queue := kafka.NewClient(kafka.DefaultOptions)
	defer queue.Close()

	//go clearExpiredAddr()

	groupID := "heart_" + fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	err := queue.ReceiveMessage(groupID, TopicName, func(i interface{}) error {

		msg := i.(kafka.Message)
		mutex.Lock()
		addrList[msg.Message] = time.Now()
		mutex.Unlock()

		return nil
	})
	if err != nil {
		return
	}

}

// clear expired addr,10 seconds no reply
func clearExpiredAddr() {
	for {
		time.Sleep(1 * time.Second)
		mutex.Lock()
		for addr, lastTime := range addrList {
			if lastTime.Add(10 * time.Second).Before(time.Now()) { //10 shoud be config
				delete(addrList, addr)
			}
		}
		mutex.Unlock()
	}
}

// get all live addr list
func GetAddrList() []string {
	mutex.Lock()
	defer mutex.Unlock()
	rs := make([]string, 0)
	for addr, _ := range addrList {
		rs = append(rs, addr)
	}
	return rs
}
