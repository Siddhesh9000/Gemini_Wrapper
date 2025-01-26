package main

import (
	"log"
	"os"

	"main/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

// Check for Vercel environment (only used in Vercel deployment)
isVercel := os.Getenv("VERCEL")
if isVercel != "" {
	// Skip local server in Vercel environment, handle via generate.go
	log.Println("Running in Vercel, using generate.go")
} else {
	// Start the local server (for development)
	server.StartServer()
}
}
