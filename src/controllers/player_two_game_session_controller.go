package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tic-tac-toe-game/src/models"
	"tic-tac-toe-game/src/services"
)

type PlayerTwoGameSessionController struct {
	Game *models.GameSession
}

type AddPlayerTwoRequest struct {
	PlayerTwoName string `json:"player2" binding:"required"`
}

func NewPlayerTwoGameSessionController() *PlayerTwoGameSessionController {
	return &PlayerTwoGameSessionController{}
}

func (bc *PlayerTwoGameSessionController) PlayerTwoGameSessionControllerHandler(c *gin.Context) {
	var req AddPlayerTwoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameSessionId := c.Param("gameSessionId")

	gameSession, err := services.AddPlayerTwoToGameSession(gameSessionId, req.PlayerTwoName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gameSession)
}
