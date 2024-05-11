package client

import (
	"context"

	"github.com/kushal0511-not/toll_calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	Client   types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endpoint,
		Client:   c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	_, err := c.Client.Aggregate(ctx, aggReq)
	if err != nil {
		return err
	}
	return nil
}

func (c *GRPCClient) GetInvoice(ctx context.Context, invReq *types.InvoiceRequest) (*types.Invoice, error) {
	resp, err := c.Client.GetInvoice(ctx, invReq)
	if err != nil {
		return nil, err
	}
	invoice := &types.Invoice{
		OBUID:         int(resp.ObuId),
		TotalAmount:   resp.TotalAmount,
		TotalDistance: resp.TotalDistance,
	}
	return invoice, err
}
