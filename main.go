package main

import (
	"log"
	"os"

	"gin-project/config"
	"gin-project/models"
	"gin-project/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	config.ConnectDB()

	// Auto-migrate database tables
	db := config.GetDB()
	err := db.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Setup all routes
	routes.SetupAllRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
