package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kushal0511-not/toll_calculator/aggregator/client"
	"github.com/kushal0511-not/toll_calculator/types"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "listen address for HTTP server")
	aggregateServiceAddr := flag.String("aggregateServiceAddr", "http://localhost:3030", "listen address for HTTP server")
	flag.Parse()
	var (
		client  = client.NewHTTPClient(*aggregateServiceAddr)
		handler = NewInvoiceHandler(client)
	)
	http.HandleFunc("/invoices", makeAPIFunc(handler.handleGetInvoice))
	slog.Info("Gatway server listening on ", "listenAddr", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(client client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: client,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	obu, ok := r.URL.Query()["obu"]
	if !ok {
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "obu is not provided"})
	}
	if len(obu) != 1 {
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "obu is not provided"})
	}
	id, err := strconv.Atoi(obu[0])
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}
	invReq := &types.InvoiceRequest{
		ObuId: int32(id),
	}

	inv, err := h.client.GetInvoice(context.Background(), invReq)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)

}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
	}
}
