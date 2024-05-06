package main

import (
	"context"

	"github.com/kushal0511-not/toll_calculator/types"
)

type GRPCServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCServer(svc Aggregator) *GRPCServer {
	return &GRPCServer{
		svc: svc,
	}
}

func (s *GRPCServer) Aggregate(ctx context.Context, r *types.AggregateRequest) (*types.None, error) {
	distance := &types.Distance{
		OBUID: int(r.ObuId),
		Value: r.Value,
		Unix:  r.Unix,
	}

	return &types.None{}, s.svc.AggregateDistance(*distance)
}
