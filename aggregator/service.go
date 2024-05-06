package main

import (
	"github.com/kushal0511-not/toll_calculator/types"
)

const basPrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}
type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}
type InvoioceAggregator struct {
	store Storer
}

func NewInvoioceAggregator(store Storer) Aggregator {
	return &InvoioceAggregator{
		store: store,
	}
}

func (i *InvoioceAggregator) AggregateDistance(distance types.Distance) error {
	return i.store.Insert(distance)
}
func (i *InvoioceAggregator) CalculateInvoice(id int) (*types.Invoice, error) {
	dist, err := i.store.Get(id)
	if err != nil {
		return nil, err
	}
	return &types.Invoice{
		OBUID:         id,
		TotalDistance: dist,
		TotalAmount:   basPrice * dist,
	}, nil
}
