package main

import (
	"math"

	"github.com/kushal0511-not/toll_calculator/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculatorService struct {
	prevPoint []float64
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{
		prevPoint: make([]float64, 2),
	}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	if len(s.prevPoint) > 0 {
		distance := euclideanDistance(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Long)
		return distance, nil
	}
	s.prevPoint = []float64{data.Lat, data.Long}
	return 0.0, nil
}

func euclideanDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
