package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"dex-analyzer/internal/api"
)

func main() {
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	// Initialize API server
	server := api.NewServer()

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/positions", server.GetPositions).Methods("GET")
	r.HandleFunc("/analyze", server.AnalyzeWithASI).Methods("GET")

	// Start server
	address := fmt.Sprintf(":%s", *port)
	log.Printf("Starting server on %s", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
