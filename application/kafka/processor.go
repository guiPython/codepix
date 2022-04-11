package kafka

import (
	"fmt"

	"github.com/jinzhu/gorm"
	ckafka "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaProcessor struct {
	Database        *gorm.DB
	Producer        *ckafka.Producer
	DeliveryChannel chan ckafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChannel chan ckafka.Event) *KafkaProcessor {
	processor := &KafkaProcessor{
		Database:        database,
		Producer:        producer,
		DeliveryChannel: deliveryChannel,
	}
	return processor
}

func (processor *KafkaProcessor) Consume() {
	var consumer *ckafka.Consumer
	var err error

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "consumergroup",
		"auto.offset.reset": "earliest",
	}

	if consumer, err = ckafka.NewConsumer(&configMap); err != nil {
		panic(err)
	}

	topics := []string{"teste"}
	consumer.SubscribeTopics(topics, nil)

	for {
		if msg, err := consumer.ReadMessage(-1); err == nil {
			processor.processMessage((msg))
		}
	}
}

func (processo *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := "transactions"
	transactionConfirmationTopic := "transaction_confirmation"

	switch topic := *&msg.TopicPartition.Topic; topic {
	case &transactionsTopic:
	case &transactionConfirmationTopic:
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (processor *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	return nil
}
