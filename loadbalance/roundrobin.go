package loadbalance

// 轮询
type RoundRobin struct {
	addrs []string

	curIdx int
}

func (r *RoundRobin) Add(param []string) error {
	if len(param) != 1 {
		return ErrParam
	}
	r.addrs = append(r.addrs, param[0])
	return nil
}

func (r *RoundRobin) Get(string) (string, error) {
	if len(r.addrs) == 0 || r.curIdx >= len(r.addrs) {
		return "", ErrNoAddr
	}
	addr := r.addrs[r.curIdx]
	r.curIdx = (r.curIdx + 1) % len(r.addrs)
	return addr, nil
}
