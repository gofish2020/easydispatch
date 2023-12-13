package loadbalance

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistHash(t *testing.T) {

	//
	chash := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data)) // 仅仅转换为整数
		return uint32(i)
	})

	chash.Add([]string{"3", "5"})
	t.Log(chash.hashSlots) //hashSlots: [3 5 13 15 23 25]

	testCase := map[string]string{
		"2":  "3", // 3 > 2
		"4":  "5", // 5 > 4
		"13": "3", // 13 >= 13
		"19": "3", // 23 > 19
		"29": "3",
	}

	for k, v := range testCase {
		actual, _ := chash.Get(k)
		assert.Equal(t, v, actual)
	}

	chash.Add([]string{"8"})
	t.Log(chash.hashSlots) //hashSlots: [3 5 8 13 15 18 23 25 28]

	testCase["27"] = "8" // 28 > 27
	for k, v := range testCase {
		actual, _ := chash.Get(k)
		assert.Equal(t, v, actual)
	}
}
