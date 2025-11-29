package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupAllRoutes initializes all route groups
func SetupAllRoutes(router *gin.Engine) {
	SetupPublicRoutes(router)
	SetupAuthRoutes(router)
	SetupProtectedRoutes(router)
}
