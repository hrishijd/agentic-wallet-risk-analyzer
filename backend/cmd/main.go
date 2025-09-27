package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	// Set up Gin router
	r := gin.Default()

	// CORS middleware for localhost
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// Endpoints
	r.GET("/positions", func(c *gin.Context) {
		server.GetPositions(c.Writer, c.Request)
	})

	r.GET("/analyze", func(c *gin.Context) {
		server.AnalyzeWithASI(c.Writer, c.Request)
	})

	// Start server
	address := fmt.Sprintf(":%s", *port)
	log.Printf("Starting server on %s", address)
	if err := r.Run(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
