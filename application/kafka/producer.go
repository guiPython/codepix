package kafka

import (
	"fmt"

	ckafka "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func NewKafkaProducer() *ckafka.Producer {
	var producer *ckafka.Producer
	var err error
	config := ckafka.ConfigMap{"bootstrap.servers": "localhost:9092"}
	if producer, err = ckafka.NewProducer(&config); err != nil {
		panic(err)
	}
	return producer
}

func Publish(msg, topic string, producer *ckafka.Producer, deliveryChannel chan ckafka.Event) error {
	message := ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic: &topic,
		},
		Value: []byte(msg),
	}

	if err := producer.Produce(&message, deliveryChannel); err != nil {
		return err
	}
	return nil
}

func DeliveryReport(deliveryChannel chan ckafka.Event) {
	for event := range deliveryChannel {
		switch ev := event.(type) {
		case *ckafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Delivery failed:", ev.TopicPartition)
			} else {
				fmt.Println("Delivered message to:", ev.TopicPartition)
			}
		}
	}
}
