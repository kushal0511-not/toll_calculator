package main

import (
	"log/slog"
	"time"

	"github.com/kushal0511-not/toll_calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsMiddleware struct {
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next           Aggregator
}

type LoggingMiddleware struct {
	next Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregator",
		Help:      "Total number of errors",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
		Help:      "Total number of errors",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregator",
		Help:      "Total number of requests",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
		Help:      "Total number of requests",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculae",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		next:           next,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}

}

func NewLoggingMiddleware(next Aggregator) Aggregator {
	return &LoggingMiddleware{
		next: next,
	}
}

func (mm *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	start := time.Now()
	defer func(start time.Time) {
		mm.reqCounterAgg.Inc()
		mm.reqLatencyAgg.Observe(time.Since(start).Seconds())
		if err != nil {
			mm.errCounterAgg.Inc()
		}
	}(start)
	err = mm.next.AggregateDistance(distance)
	return

}

func (mm *MetricsMiddleware) CalculateInvoice(id int) (inv *types.Invoice, err error) {
	start := time.Now()
	defer func(start time.Time) {
		mm.reqCounterCalc.Inc()
		mm.reqLatencyCalc.Observe(time.Since(start).Seconds())
		if err != nil {
			mm.errCounterCalc.Inc()
		}
	}(start)
	inv, err = mm.next.CalculateInvoice(id)
	return
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
		slog.Info("Calculate Invoice Logging Middleware", "id", id, "totalDistance", distance, "totalAmount", amount, "Took", time.Since(start).String())
	}(start)
	inv, err = lm.next.CalculateInvoice(id)
	return
}
