package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Message struct {
	Topic   string
	Key     string
	Message string
}

func NewClient(option Options) *Client {
	return &Client{
		option: option,
	}
}

type Client struct {
	mu       sync.Mutex
	producer sarama.SyncProducer
	option   Options
}

func (c *Client) SendMessage(msg Message) (int32, int64, error) {

	c.newProducer()
	message := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.StringEncoder(msg.Key),
		Value: sarama.StringEncoder(msg.Message),
	}
	return c.producer.SendMessage(message)

}

func (c *Client) Close() {
	if c.producer != nil {
		c.producer.Close()
	}
}

func (c *Client) newProducer() error {

	if c.producer != nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.producer != nil {
		return nil
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = c.option.RequiredAcks
	config.Producer.Retry.Max = c.option.MaxRetry
	config.Producer.Timeout = c.option.Timeout
	config.Producer.Return.Successes = true

	// 创建生产者
	producer, err := sarama.NewSyncProducer(c.option.Brokers, config)
	if err != nil {
		return err
	}

	c.producer = producer
	return nil
}

func (c *Client) ReceiveMessage(groupID, topic string, f func(interface{}) error) error {

	if f == nil {
		return errors.New("f is nil")
	}
	config := sarama.NewConfig()
	config.Version = c.option.Version // 设置 Kafka 版本

	switch c.option.Rebalance {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()} // rebalance sticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()} // rebalance sticky
	}

	switch c.option.OffsetInitial {
	case sarama.OffsetNewest:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	case sarama.OffsetOldest:
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	default:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest

	}

	consumer := Consumer{
		ready: make(chan bool),
		f:     f,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := sarama.NewConsumerGroup(c.option.Brokers, groupID, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("Error closing consumer group: %v\n", err)
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{topic}, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	log.Println("Sarama consumer up and running!...")
	wg.Wait()
	return nil
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	f     func(interface{}) error
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}

			msg := Message{
				Topic:   message.Topic,
				Key:     string(message.Key),
				Message: string(message.Value),
			}
			err := consumer.f(msg)
			if err != nil {
				continue
			}

			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
