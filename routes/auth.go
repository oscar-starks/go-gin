package routes

import (
	"gin-project/handlers"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures authentication-related routes
func SetupAuthRoutes(router *gin.Engine) {
	// Authentication routes group
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register/", handlers.Register)
		authRoutes.POST("/login/", handlers.Login)
	}
}
