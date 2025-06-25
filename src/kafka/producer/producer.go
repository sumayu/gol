package main

import (
	"fmt"
	"main/src/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	config:= &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", 
	}
	producer,err  := kafka.NewProducer(config)
	if err != nil {
		logger.Logger.Error("kafka INIT ERROR", err)
	}
	defer producer.Close()
	topic := "test-topic"
	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("сообщение%d",i)
err := producer.Produce(&kafka.Message{ TopicPartition: kafka.TopicPartition{
	Topic: &topic, 
	Partition: kafka.PartitionAny,

}}, nil)
if err != nil {
	logger.Logger.Error("Producer cant push message to kafka PRODUCER.ERROR", err)
} else { 
	fmt.Println("message push", message)
}
producer.Flush(15*1000)

	}
}