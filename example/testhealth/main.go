package main

import (
	"fmt"
	"time"

	"github.com/gofish2020/easydispatch/health"
)

func main() {
	health.AddAddr("https://www.sina.com.cn/", "https://www.baidu.com/", "http://www.aajklsdfjklsd")
	go health.HealthCheck()

	time.Sleep(50 * time.Second)

	alist := health.GetAliveAddrList()
	for i := 0; i < len(alist); i++ {
		fmt.Println(alist[i])
	}

	var block = make(chan bool)
	<-block
}
