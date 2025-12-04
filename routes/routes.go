package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupAllRoutes(router *gin.Engine) {
	SetupPublicRoutes(router)
	SetupAuthRoutes(router)
	SetupProtectedRoutes(router)
	SetupChatRoutes(router)
	SetupNotificationRoutes(router)
	SetupWebSocketRoutes(router)
}
