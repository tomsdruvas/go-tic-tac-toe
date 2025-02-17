package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"src/src/controllers"
)

func main() {
	healthCheckController := controllers.NewHealthCheckController()
	boardController := controllers.NewCreateGameController()
	getBoardController := controllers.NewGetGameSessionController()

	router := gin.Default()

	router.GET("/health", healthCheckController.HealthCheckHandler)
	router.POST("/game-session", boardController.CreateGameControllerHandler)
	router.GET("/game-session/:gameSessionId", getBoardController.GetGameSessionControllerHandler)

	router.Run(":8080")
}
