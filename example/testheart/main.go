package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofish2020/easydispatch/heart"
	"github.com/gofish2020/easydispatch/loadbalance"
)

// 利用消息队列作为服务存活监测手段
func main() {

	// 服务负责发送心跳
	go heart.RunHeartBeat()
	// 调度器负责接收心跳
	go heart.ListenHeartbeat()

	// 利用负载均衡获取

	lb := loadbalance.LoadBalanceFactory(loadbalance.BalanceTypeRand)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			for _, v := range heart.GetAddrList() {
				lb.Add([]string{v})
			}
			fmt.Println(lb.Get(""))
		}

	}()

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGTERM)
	<-sigusr1
}
