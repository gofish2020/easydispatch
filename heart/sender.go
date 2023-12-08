package heart

import (
	"errors"
	"net"
	"time"

	"github.com/gofish2020/easydispatch/heart/kafka"
)

const TopicName = "heart-beat"

/*
	负责往kafka中发送消息
*/

func RunHeartBeat() {
	client := kafka.NewClient(kafka.DefaultOptions)
	addr, _ := ExternalIP()
	for {
		client.SendMessage(kafka.Message{
			Topic:   TopicName,
			Key:     "",
			Message: addr.String(),
		})
		// 每秒发送ip信息到Kafka中
		time.Sleep(1 * time.Second)
	}

}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
