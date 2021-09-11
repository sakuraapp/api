package main

import (
	"github.com/joho/godotenv"
	"github.com/sakuraapp/api/server"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	server.Start(port)
}