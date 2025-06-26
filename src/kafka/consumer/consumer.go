package main

import (
	"main/src/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)
func main()  {
logger.Init()
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id": "mygroup",
		"auto.offset.reset": "earliest", 
	}

consumer, err := kafka.NewConsumer(config)
if err!= nil {
	logger.Logger.Error("InitConsumerErrorKAFKA", err)
} else {
	logger.Logger.Debug("ConsumerInitKAFKA")
}
defer consumer.Close()
topic := "test-topic"
err = consumer.SubscribeTopics([]string{topic}, nil)
if err != nil {
	    logger.Logger.Error("Ошибка подписки на топик", "topic", topic, "error", err)

}
for {
	msg, err := consumer.ReadMessage(-1)
	if err != nil {
		logger.Logger.Error("Ошибка чтения", err)

	}
	if msg == nil { 
    logger.Logger.Warn("Получено nil-сообщение")
    continue
}
	logger.Logger.Debug("Получено сообщение", string(msg.Value)) //TO-DO решить проблему с выводом данных в msg.Value-> почему оно не сохраняет сообщения
}

}
