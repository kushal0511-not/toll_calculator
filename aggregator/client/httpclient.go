package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kushal0511-not/toll_calculator/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	httpc := http.DefaultClient
	body, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("aggregate Invoice Error %v", err)
	}
	return nil
}

func (c *HTTPClient) GetInvoice(ctx context.Context, invReq *types.InvoiceRequest) (*types.Invoice, error) {
	httpc := http.DefaultClient
	body, err := json.Marshal(invReq)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s?obu=%d", c.Endpoint, "invoices", invReq.ObuId)
	slog.Info(
		"GetInvoice",
		slog.String("endpoint", c.Endpoint),
		slog.String("id", string(invReq.ObuId)),
	)
	slog.Info(endpoint)
	req, err := http.NewRequest("GET", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get Invoice Error %v", err)
	}
	var invoice types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, err
	}
	return &invoice, nil

}
