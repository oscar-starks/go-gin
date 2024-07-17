package main

import (
	// "github.com/gin-gonic/gin"
	"github.com/oscar-starks/go-gin/routes"
)

func main() {
	// r := gin.Default()
	router := routes.SetupRouter()
	
	router.Run("localhost:8000") // listen and serve on 0.0.0.0:8080
}
