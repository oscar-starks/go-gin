package main

import (
	"log"
	"os"

	"gin-project/config"
	"gin-project/handlers"
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
	config.ConnectRedis()

	// Initialize WebSocket connection manager
	handlers.InitializeWebSocketManager()

	// Auto-migrate database tables
	db := config.GetDB()
	err := db.AutoMigrate(&models.User{}, &models.Room{}, &models.ChatRequest{}, &models.Notification{}, &models.Message{})

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

	// // Clear all room_members associations first
	// result := config.DB.Exec("DELETE FROM room_members")
	// if result.Error != nil {
	// 	log.Fatal("Failed to clear room members:", result.Error)
	// }
	// log.Printf("Cleared %d room members from the database.\n", result.RowsAffected)

	// // Clear all chat requests
	// result = config.DB.Unscoped().Delete(&models.ChatRequest{}, "1 = 1")
	// if result.Error != nil {
	// 	log.Fatal("Failed to clear chat requests:", result.Error)
	// }
	// log.Printf("Cleared %d chat requests from the database.\n", result.RowsAffected)

	// // Clear all rooms
	// result = config.DB.Unscoped().Delete(&models.Room{}, "1 = 1")
	// if result.Error != nil {
	// 	log.Fatal("Failed to clear rooms:", result.Error)
	// }
	// log.Printf("Cleared %d rooms from the database.\n", result.RowsAffected)

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
