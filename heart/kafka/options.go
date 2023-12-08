package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

type Options struct {
	Brokers []string

	// producer
	RequiredAcks sarama.RequiredAcks // NoResponse WaitForLocal

	MaxRetry int
	Timeout  time.Duration

	// consumer
	Version       sarama.KafkaVersion
	Rebalance     string // sticky roundrobin range
	OffsetInitial int64  // OffsetNewest OffsetOldest

}

var DefaultOptions = Options{
	Brokers:       []string{"kafka:9092"},
	RequiredAcks:  sarama.WaitForLocal,
	MaxRetry:      3,
	Timeout:       5 * time.Second,
	Version:       sarama.V2_8_1_0,
	Rebalance:     "sticky",
	OffsetInitial: sarama.OffsetNewest,
}
