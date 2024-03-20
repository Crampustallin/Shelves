package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/Crampustallin/Shelves/internal/handlers"
)

func main() {
	if len(os.Args) <= 1 {
		log.Println("Not enough arguments\n The usage is main.go [order nums]")
		return 
	}
	args := os.Args[1:]

	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	handlers.HandleOrdersQuery(args)
}
