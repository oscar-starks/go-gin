package routes

import (
	"gin-project/handlers"
	"gin-project/middleware"

	"github.com/gin-gonic/gin"
)

func SetupWebSocketRoutes(router *gin.Engine) {
	router.GET("/ws", middleware.WebSocketAuthMiddleware(), handlers.HandleWebSocket)
}
