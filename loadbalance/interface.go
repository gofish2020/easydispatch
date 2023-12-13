package loadbalance

// 负载均衡器
type LoadBalance interface {
	Add([]string) error
	Get(string) (string, error)
}

type BalanceType = int

const (
	BalanceTypeRand        BalanceType = iota // 随机
	BalanceTypeRR                             // 轮询
	BalanceTypeWeigthRR                       // 加权轮询
	BalanceTypeConsistHash                    // 一致性hash

)

func LoadBalanceFactory(bType BalanceType) LoadBalance {
	switch bType {
	case BalanceTypeRand:
		return &Rand{}
	case BalanceTypeRR:
		return &RoundRobin{}
	case BalanceTypeWeigthRR:
		return &WeigthRoundRobin{}
	case BalanceTypeConsistHash:
		return New(50, nil)
	}
	return &Rand{}
}
