package main

import (
	"log"

	"api-server/internal/database"
	"api-server/internal/handlers"
	"api-server/pkg/config"
)

func main() {
	config.LoadConfig()

	database.InitDB()
	defer database.CloseDB()

	router := handlers.SetupRouter()

	log.Println("Server is running on port 8080")
	log.Fatal(router.Run(":8080"))
}
