package main

import (
	"log"
	"log/slog"

	"github.com/kushal0511-not/toll_calculator/aggregator/client"
)

type DistanceCalculator struct {
	consumer DataConsumer
}

var (
	topic    = "OBUData"
	endpoint = "http://127.0.0.1:3030"
	// grpcEndpoint = ":3031"
)

func main() {
	var (
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLoggingMiddleware(svc)
	httpClient := client.NewHTTPClient(endpoint)
	// grpcClient, err := client.NewGRPCClient(grpcEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	dc := NewDistanceCalculator(svc, httpClient)
	if err := dc.consumer.ConsumeData(); err != nil {
		slog.Error("Failed to consume", err)
	}
}

func NewDistanceCalculator(svc CalculatorServicer, client client.Client) *DistanceCalculator {
	c, err := NewKafkaConsumer(topic, svc, client)
	if err != nil {
		log.Fatal(err)
	}
	return &DistanceCalculator{
		consumer: c,
	}
}
