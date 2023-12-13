package loadbalance

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 哈希函数
type Hasher func(data []byte) uint32

type UInt32Slice []uint32

func (u UInt32Slice) Len() int {
	return len(u)
}

func (u UInt32Slice) Less(i, j int) bool {
	return u[i] < u[j]
}

func (u UInt32Slice) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type ConsistHash struct {
	hashSlots UInt32Slice       // 哈希环
	mapped    map[uint32]string // 哈希值映射实际值
	hasher    Hasher

	replicas int
}

// replicas:副本数量
func New(replicas int, hasher Hasher) *ConsistHash {
	consistHash := &ConsistHash{
		mapped:   make(map[uint32]string),
		replicas: replicas,
		hasher:   hasher,
	}
	if consistHash.hasher == nil {
		consistHash.hasher = crc32.ChecksumIEEE
	}
	return consistHash
}

func (c *ConsistHash) Add(addrs []string) error {
	for _, addr := range addrs {
		for i := 0; i < c.replicas; i++ { // 为了addr散落的更均匀（避免过度倾斜）
			// 计算hash
			hashSlot := c.hasher([]byte(strconv.Itoa(i) + addr))
			// 记录hash槽值
			c.hashSlots = append(c.hashSlots, hashSlot)
			// 映射槽值和实际值
			c.mapped[hashSlot] = addr
		}
	}
	// 排序hash环
	sort.Sort(c.hashSlots)
	return nil
}

func (c *ConsistHash) Get(key string) (string, error) {
	hashSlot := c.hasher([]byte(key))
	//二分搜索，找到大于等于 hashSlot的第一个元素索引
	idx := sort.Search(len(c.hashSlots), func(i int) bool {
		return c.hashSlots[i] >= hashSlot
	})
	// 越界，环就从头开始
	if idx == len(c.hashSlots) {
		idx = 0
	}
	return c.mapped[c.hashSlots[idx]], nil
}
