package main

import (
	"log/slog"
	"time"

	"github.com/kushal0511-not/toll_calculator/types"
)

type LoggingMiddleware struct {
	next DataProducer
}

func NewLoggingMiddleware(next DataProducer) *LoggingMiddleware {
	return &LoggingMiddleware{
		next: next,
	}
}

func (lm *LoggingMiddleware) ProduceData(data types.OBUData) error {
	start := time.Now()
	defer func(start time.Time) {
		slog.Info("Prducer Logging Middleware",
			"OBUId", data.OBUID,
			"Lat", data.Lat,
			"Long", data.Long,
			"Took", time.Since(start).String(),
		)
	}(start)

	return lm.next.ProduceData(data)
}
