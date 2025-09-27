package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"dex-analyzer/internal/api"
)

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := flag.String("port", getEnvOrDefault("PORT", "8080"), "Port to run the server on")
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
