package loadbalance

import (
	"math"
	"strconv"
)

type WeigthRoundRobin struct {
	weightAddrs []*weightAddr
}

type weightAddr struct {
	addr      string // 地址
	weight    int    // 权重
	curWeight int    // 计算使用
}

func (w *WeigthRoundRobin) Add(param []string) error {
	if len(param) != 2 {
		return ErrParam
	}

	weight, err := strconv.Atoi(param[1])
	if err != nil {
		return err
	}
	w.weightAddrs = append(w.weightAddrs, &weightAddr{
		addr:      param[0],
		weight:    weight,
		curWeight: 0,
	})

	return nil
}

func (w *WeigthRoundRobin) Get(string) (string, error) {

	if len(w.weightAddrs) == 0 {
		return "", ErrNoAddr
	}

	maxWeight := math.MinInt
	idx := 0
	sumWeight := 0 // 权重总和
	for k, weightAddr := range w.weightAddrs {

		sumWeight += weightAddr.weight // 权重总和

		weightAddr.curWeight += weightAddr.weight // 加上权重
		if weightAddr.curWeight > maxWeight {     // 记录本次最大值
			maxWeight = weightAddr.curWeight
			idx = k
		}
	}

	w.weightAddrs[idx].curWeight -= sumWeight // 减去权重总和
	return w.weightAddrs[idx].addr, nil       // 返回最大权重的结果
}
