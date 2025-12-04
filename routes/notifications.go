package routes

import (
	"gin-project/handlers"
	"gin-project/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(router *gin.Engine) {
	protectedRoutes := router.Group("/notifications")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.GET("/", handlers.ListNotifications)
		protectedRoutes.POST("/mark_read/:notificationID/", handlers.MarkNotificationAsRead)
	}
}
