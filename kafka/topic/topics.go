package topic

import (
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

type Topic struct{}

func (t *Topic) CreateTopic(topicName string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function", r)
		}
	}()
	conn, err := kafka.Dial("tcp", "localhost:19092") // Dialing the Kafka Cluster
	if err != nil {
		return false
	}
	controller, err := conn.Controller() // Getting a Kafka Bootstrap Broker
	if err != nil {
		return false
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))) // Dialing the Kafka Broker
	if err != nil {
		return false
	}
	kafkaTopic := []kafka.TopicConfig{
		{
			Topic:             topicName,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(kafkaTopic...)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
