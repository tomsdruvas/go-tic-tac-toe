package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"src/src/controllers"
)

func main() {
	healthCheckController := controllers.NewHealthCheckController()
	boardController := controllers.NewCreateGameController()

	router := gin.Default()

	router.GET("/health", healthCheckController.HealthCheckHandler)
	router.GET("/hello", boardController.CreateGameControllerHandler)

	router.Run(":8080")
}
