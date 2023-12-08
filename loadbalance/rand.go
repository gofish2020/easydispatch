package loadbalance

import (
	"math/rand"
	"time"
)

// 随机
type Rand struct {
	addrs []string
}

func (r *Rand) Add(param []string) error {

	if len(param) != 1 {
		return ErrParam
	}
	r.addrs = append(r.addrs, param[0])
	return nil
}

func (r *Rand) Get() (string, error) {
	if len(r.addrs) == 0 {
		return "", ErrNoAddr
	}
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(r.addrs))
	return r.addrs[idx], nil
}
