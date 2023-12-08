package main

import (
	"fmt"

	"github.com/gofish2020/easydispatch/loadbalance"
)

func main() {
	l1 := loadbalance.LoadBalanceFactory(loadbalance.BalanceTypeRand)
	l1.Add([]string{"127.0.0.1:80"})
	l1.Add([]string{"127.0.0.1:81"})
	l1.Add([]string{"127.0.0.1:82"})

	fmt.Println("random get")
	for i := 0; i < 10; i++ { // 随机出现
		fmt.Println(l1.Get())
	}

	l2 := loadbalance.LoadBalanceFactory(loadbalance.BalanceTypeRR)
	l2.Add([]string{"127.0.0.1:80"})
	l2.Add([]string{"127.0.0.1:81"})
	l2.Add([]string{"127.0.0.1:82"})
	fmt.Println("RoundRobin get")
	for i := 0; i < 10; i++ { // 轮流出现
		fmt.Println(l2.Get())
	}

	l3 := loadbalance.LoadBalanceFactory(loadbalance.BalanceTypeWeigthRR)
	l3.Add([]string{"127.0.0.1:80", "2"})
	l3.Add([]string{"127.0.0.1:81", "2"})
	l3.Add([]string{"127.0.0.1:82", "6"})
	fmt.Println("Weight RoundRobin get")

	count := make(map[string]int)
	for i := 0; i < 200000; i++ {
		result, _ := l3.Get()
		count[result]++
	}
	fmt.Println(count) // map[127.0.0.1:80:40000 127.0.0.1:81:40000 127.0.0.1:82:120000]
}
