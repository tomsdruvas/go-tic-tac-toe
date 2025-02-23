package main

import (
	"github.com/gin-gonic/gin"
	"log"
	_ "net/http"
	"tic-tac-toe-game/src/controllers"
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

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
