package main

import (
	"github.com/oscar-starks/go-gin/routes"
)

func main() {
	router := routes.SetupRouter()
	
	router.Run("localhost:8000") // listen and serve on 0.0.0.0:8000
}
