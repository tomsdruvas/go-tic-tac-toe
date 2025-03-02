package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"tic-tac-toe-game/src/container_setup"
	"tic-tac-toe-game/src/controllers"
	"tic-tac-toe-game/src/database"
	"tic-tac-toe-game/src/models"
	"tic-tac-toe-game/src/websockets"
)

func main() {

	container := container_setup.BuildContainer()
	if container == nil {
		log.Fatal("Failed to build dependency injection container")
	}

	err := container.Invoke(func(
		healthCheckController *controllers.HealthCheckController,
		gameSessionController *controllers.GameSessionController,
		webSocketService *websockets.WebSocketService,
		db *database.InMemoryGameSessionDB) {

		db.StoreSession(models.GameSession{
			SessionId: "123",
		})

		router := gin.Default()
		webSocketService.StartWebSocketServer()

		router.GET("/health", healthCheckController.HealthCheckHandler)
		router.POST("/game-session", gameSessionController.CreateGameSessionHandler)
		router.GET("/game-session/:gameSessionId", gameSessionController.GetGameSessionHandler)
		router.POST("/game-session/:gameSessionId/players", gameSessionController.PlayerTwoGameSessionHandler)
		router.POST("/game-session/:gameSessionId/move", gameSessionController.SubmitMoveHandler)

		err := router.Run(":8080")
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

}
