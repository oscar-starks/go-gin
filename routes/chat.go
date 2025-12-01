package routes

import (
	"gin-project/handlers"
	"gin-project/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(router *gin.Engine) {
	protectedRoutes := router.Group("/chat")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.GET("/has_room/", handlers.HasRoom)
	}
}
