package main

import (
	"os"

	"github.com/doeta/Kasir-Go/configs"
	_ "github.com/doeta/Kasir-Go/docs"
	"github.com/doeta/Kasir-Go/routes"
	"github.com/joho/godotenv"
)


func main() {
    // Load .env untuk local development
    godotenv.Load()

    // Konek Database
    configs.ConnectDB()
    
    // Setup Router
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}