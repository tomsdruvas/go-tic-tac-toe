package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tic-tac-toe-game/src/services"
)

type CreateGameController struct {
}

type CreateGameRequest struct {
	PlayerName string `json:"player1" binding:"required"`
}

func NewCreateGameController() *CreateGameController {
	return &CreateGameController{}
}

func (bc *CreateGameController) CreateGameControllerHandler(c *gin.Context) {
	var req CreateGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := services.CreateTicTacToeGameSession(req.PlayerName)

	c.JSON(http.StatusCreated, session)
}
