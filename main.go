package main

import (
	"log"

	"github.com/joho/godotenv"
	"main/server" 
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server.StartServer()
}
