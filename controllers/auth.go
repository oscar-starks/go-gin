package auth


import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func GetHome(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "This is the homepage endpoint!",
    })
}

func Ping(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "This is a ping test",
    })
}
