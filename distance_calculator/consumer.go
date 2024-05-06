package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/kushal0511-not/toll_calculator/aggregator/client"
	"github.com/kushal0511-not/toll_calculator/types"
)

type DataConsumer interface {
	ConsumeData() error
}

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (DataConsumer, error) {
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
		aggClient:   aggClient,
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
		req := &types.AggregateRequest{
			ObuId: int32(data.OBUID),
			Value: distance,
			Unix:  time.Now().Unix(),
		}
		if err := c.aggClient.Aggregate(context.Background(), req); err != nil {
			slog.Error("Error sending data to aggregator")
			continue
		}
	}

	c.consumer.Close()
	return nil
}
