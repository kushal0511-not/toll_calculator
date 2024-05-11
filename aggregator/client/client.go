package client

import (
	"context"

	"github.com/kushal0511-not/toll_calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, *types.InvoiceRequest) (*types.Invoice, error)
}
