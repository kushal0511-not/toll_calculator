package main

import (
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/kushal0511-not/toll_calculator/types"
	"google.golang.org/grpc"
)

func main() {
	httplistenAddr := flag.String("httplistenaddr", ":3030", "HTTP listen address of http server")
	grpclistenAddr := flag.String("grpclistenaddr", ":3031", "GRPC listen address of http server")

	flag.Parse()

	store := NewMemoryStore()
	var (
		svc = NewInvoioceAggregator(store)
	)
	svc = NewLoggingMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(svc, grpclistenAddr))
	}()

	log.Fatal(makeHTTPTransport(svc, httplistenAddr))

}

func makeHTTPTransport(svc Aggregator, listenAddr *string) error {
	slog.Info("HTTP Server:", "Addr", *listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoices", handleGetInvoice(svc))
	return http.ListenAndServe(*listenAddr, nil)
}

func makeGRPCTransport(svc Aggregator, listenAddr *string) error {
	slog.Info("GRPC Server:", "Addr", *listenAddr)
	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	//Make GRPC server with default options
	server := grpc.NewServer([]grpc.ServerOption{}...)
	//Register our service
	types.RegisterAggregatorServer(server, NewGRPCServer(svc))
	return server.Serve(ln)

}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "obu is not provided"})
			return
		}
		obuID := value[0]
		id, err := strconv.Atoi(obuID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			return
		}
		invoice, err := svc.CalculateInvoice(id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)

	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
