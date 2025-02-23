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
	playerTwoGameSessionController := controllers.NewPlayerTwoGameSessionController()
	submitMoveController := controllers.NewSubmitMoveController()

	router := gin.Default()

	router.GET("/health", healthCheckController.HealthCheckHandler)
	router.POST("/game-session", boardController.CreateGameControllerHandler)
	router.GET("/game-session/:gameSessionId", getBoardController.GetGameSessionControllerHandler)
	router.POST("/game-session/:gameSessionId/players", playerTwoGameSessionController.PlayerTwoGameSessionControllerHandler)
	router.POST("/game-session/:gameSessionId", submitMoveController.SubmitMoveControllerHandler)

	router.Run(":8080")
}
