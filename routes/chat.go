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
		protectedRoutes.POST("/request/:userID/", handlers.RequestChat)
		protectedRoutes.GET("/received_requests/", handlers.ListReceivedChatRequests)
		protectedRoutes.GET("/sent_requests/", handlers.ListSentChatRequests)
	}
}
