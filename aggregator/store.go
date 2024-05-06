package main

import (
	"fmt"
	"log/slog"

	"github.com/kushal0511-not/toll_calculator/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	slog.Info("Processing and inserting distance ", "distance", distance)
	m.data[distance.OBUID] += distance.Value
	return nil
}

func (m *MemoryStore) Get(obuID int) (float64, error) {

	value, ok := m.data[obuID]
	if !ok {
		return 0.0, fmt.Errorf("OBUID not found %d", obuID)
	}
	return value, nil
}
