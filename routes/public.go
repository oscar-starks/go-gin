package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupPublicRoutes(router *gin.Engine) {
	// Basic route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Gin + GORM + PostgreSQL API!",
		})
	})

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})
}
