package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/kushal0511-not/toll_calculator/types"
)

type DataConsumer interface {
	ConsumeData() error
}

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, svc CalculatorServicer) (DataConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Println("Error creating Kafka")
		panic(err)
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		isRunning:   true,
		calcService: svc,
	}, nil
}

func (c *KafkaConsumer) ConsumeData() error {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			// fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			return err
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			slog.Error("Error unmarshalling")
			continue
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			slog.Error("Error calculating distance")
			continue
		}
		fmt.Printf("Calculated Distance: %.2f\n", distance)
	}

	c.consumer.Close()
	return nil
}
