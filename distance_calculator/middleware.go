package main

import (
	"log/slog"
	"time"

	"github.com/kushal0511-not/toll_calculator/types"
)

type LoggingMiddleware struct {
	next CalculatorServicer
}

func NewLoggingMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LoggingMiddleware{
		next: next,
	}
}

func (lm *LoggingMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	start := time.Now()
	defer func(start time.Time) {
		slog.Info("Calculate Distance Logging Middleware",
			"OBUId", data.OBUID,
			"Distance", dist,
			"Took", time.Since(start).String(),
		)
	}(start)
	dist, err = lm.next.CalculateDistance(data)
	return
}
