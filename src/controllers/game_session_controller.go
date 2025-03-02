package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tic-tac-toe-game/src/services"
	"tic-tac-toe-game/src/websockets"
)

type GameSessionController struct {
	gameSessionService *services.GameSessionService
	webSocketService   *websockets.WebSocketService
}

type CreateGameRequest struct {
	PlayerName string `json:"player1" binding:"required"`
}

func NewGameSessionController(
	gameSessionService *services.GameSessionService,
	webSocketService *websockets.WebSocketService,
) *GameSessionController {
	return &GameSessionController{
		gameSessionService: gameSessionService,
		webSocketService:   webSocketService,
	}
}

func (controller *GameSessionController) CreateGameSessionHandler(c *gin.Context) {
	var req CreateGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := controller.gameSessionService.CreateTicTacToeGameSession(req.PlayerName)

	c.JSON(http.StatusCreated, session)
}

func (controller *GameSessionController) GetGameSessionHandler(c *gin.Context) {
	gameSessionId := c.Param("gameSessionId")
	gameSession, err := controller.gameSessionService.RetrieveTicTacToeGameSession(gameSessionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game session not found"})
		return
	}

	c.JSON(http.StatusOK, gameSession)
}

type AddPlayerTwoRequest struct {
	PlayerTwoName string `json:"player2" binding:"required"`
}

func (controller *GameSessionController) PlayerTwoGameSessionHandler(c *gin.Context) {
	var req AddPlayerTwoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameSessionId := c.Param("gameSessionId")

	gameSession, err := controller.gameSessionService.AddPlayerTwoToGameSession(gameSessionId, req.PlayerTwoName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(gameSession)

	err = (*websockets.GameSessionConnectionStore).SendMessageToGameSession(controller.webSocketService.GameSessionConnectionStore, gameSessionId, string(jsonData))
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gameSession)
}

type SubmitMoveRequest struct {
	PlayerName string `json:"playerName" binding:"required"`
	XAxis      *int   `json:"xAxis" binding:"required"`
	YAxis      *int   `json:"yAxis" binding:"required"`
}

func (controller *GameSessionController) SubmitMoveHandler(c *gin.Context) {
	var req SubmitMoveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameSessionId := c.Param("gameSessionId")

	gameSession, err := controller.gameSessionService.AddMoveToGameSession(gameSessionId, req.PlayerName, *req.XAxis, *req.YAxis)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(gameSession)

	err = (*websockets.GameSessionConnectionStore).SendMessageToGameSession(controller.webSocketService.GameSessionConnectionStore, gameSessionId, string(jsonData))
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusCreated, gameSession)
}
