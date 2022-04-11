/*
Copyright © 20	22 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/guiPython/codepix/application/kafka"
	"github.com/guiPython/codepix/infrastructure/database"
	"github.com/spf13/cobra"
	ckafka "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start consuming transactions using Apache Kafka",

	Run: func(cmd *cobra.Command, args []string) {
		producer := kafka.NewKafkaProducer()
		delivery := make(chan ckafka.Event)
		kafka.Publish("Olá kafka", "teste", producer, delivery)
		go kafka.DeliveryReport(delivery)

		database := database.ConnectDB(os.Getenv("env"))
		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, delivery)
		kafkaProcessor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)
}
