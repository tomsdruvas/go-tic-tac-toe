package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"src/src/controllers"
)

func main() {
	boardController := controllers.NewCreateGameController()

	router := gin.Default()

	router.GET("/hello", boardController.CreateGameControllerHandler)

	router.Run(":8080")
}
