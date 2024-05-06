package main

import (
	"log/slog"
	"time"

	"github.com/kushal0511-not/toll_calculator/types"
)

type LoggingMiddleware struct {
	next Aggregator
}

func NewLoggingMiddleware(next Aggregator) Aggregator {
	return &LoggingMiddleware{
		next: next,
	}
}

func (lm *LoggingMiddleware) AggregateDistance(distance types.Distance) error {
	start := time.Now()
	defer func(start time.Time) {
		slog.Info("Aggregate Distance Logging Middleware", "distance", distance, "Took", time.Since(start).String())
	}(start)
	return lm.next.AggregateDistance(distance)
}
func (lm *LoggingMiddleware) CalculateInvoice(id int) (inv *types.Invoice, err error) {
	start := time.Now()
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}
		slog.Info("Calculate Invoice Logging Middleware", "id", inv.OBUID, "totalDistance", distance, "totalAmount", amount, "Took", time.Since(start).String())
	}(start)
	inv, err = lm.next.CalculateInvoice(id)
	return
}
