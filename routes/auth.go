package routes

import (
    "github.com/oscar-starks/go-gin/controllers"
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    router.GET("/", auth.GetHome)
    router.GET("/ping", auth.Ping)

    return router
}
