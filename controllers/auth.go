package auth

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func GetHome(ctx *gin.Context) {
    ctx.JSON(http.StatusOK, gin.H{
        "message": "This is the homepage endpoint!",
    })
}

func Ping(ctx *gin.Context) {
    ctx.JSON(http.StatusOK, gin.H{
        "message": "This is a ping test",
    })
}
