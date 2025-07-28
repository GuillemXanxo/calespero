package main

import (
	"log"
	"os"

	"calespero/cmd/api"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Check required environment variables
	requiredEnvVars := []string{"DATABASE_URL", "JWT_SECRET_KEY"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("%s environment variable is not set", envVar)
		}
	}

	// Start the server
	api.Run()
}
