package handler

import (
    "fmt"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "main/server"
)

// Handler function that Vercel will call
func Handler(w http.ResponseWriter, r *http.Request) {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        http.Error(w, "Error loading .env file", http.StatusInternalServerError)
        return
    }

    // Check for Vercel environment (only used in Vercel deployment)
    isVercel := os.Getenv("VERCEL")
    if isVercel != "" {
        // Skip local server in Vercel environment, handle via generate.go
        fmt.Println("Running in Vercel, using generate.go")
        server.GenerateHandler(w, r) // Directly using generate handler for requests
    } else {
        // Start the local server (for development)
        server.StartServer()
    }
}
