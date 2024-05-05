package main

import (
	"log"
	"log/slog"
)

type DistanceCalculator struct {
	consumer DataConsumer
}

var topic = "OBUData"

func main() {
	var (
		svc = NewCalculatorService()
	)
	dc := NewDistanceCalculator(svc)
	if err := dc.consumer.ConsumeData(); err != nil {
		slog.Error("Failed to consume", err)
	}
}

func NewDistanceCalculator(svc *CalculatorService) *DistanceCalculator {
	c, err := NewKafkaConsumer(topic, svc)
	if err != nil {
		log.Fatal(err)
	}
	return &DistanceCalculator{
		consumer: c,
	}
}
