package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofish2020/easydispatch/heart/kafka"
)

/*
启动docker
docker-compose -f docker-compose-env.yml   up -d zookeeper
docker-compose -f docker-compose-env.yml   up -d kafka
*/
func main() {
	/*
		kafka命令 *.sh位于 /opt/kafka/bin目录
		./kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 3 --topic easy_topic
		./kafka-topics.sh --list --zookeeper zookeeper:2181
		./kafka-console-consumer.sh --bootstrap-server kafka:9092 --from-beginning --topic easy_topic
		./kafka-console-consumer.sh --bootstrap-server kafka:9092 --consumer-property group.id=testGroup --topic easy_topic
		./kafka-console-consumer.sh --bootstrap-server kafka:9092 --topic easy_topic
		./kafka-console-producer.sh --broker-list kafka:9092 --topic easy_topic
	*/
	Topic := "easy_topic"
	client := kafka.NewClient(kafka.DefaultOptions)

	go func() {
		time.Sleep(5 * time.Second)
		for i := 0; i < 10; i++ {
			client.SendMessage(kafka.Message{
				Topic:   Topic,
				Key:     "11",
				Message: "9999",
			})
		}

	}()

	groupID := strconv.Itoa(int(time.Now().UTC().UnixNano()))
	err1 := client.ReceiveMessage("easy_"+groupID, Topic, func(msg interface{}) error {

		message := msg.(kafka.Message)
		log.Printf("Message claimed: value = %s, key = %v, topic = %s", string(message.Message), message.Key, message.Topic)
		return nil
	})
	fmt.Println(err1)
}
