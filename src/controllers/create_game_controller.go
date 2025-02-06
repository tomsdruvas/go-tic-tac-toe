package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"src/src/models"
)

type CreateGameController struct {
	Game *models.GameSession
}

func NewCreateGameController() *CreateGameController {
	game := models.NewGameSession("tom")
	controller := &CreateGameController{Game: game}
	return controller
}

func (bc *CreateGameController) CreateGameControllerHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"player1":        bc.Game.Player1,
		"player2":        bc.Game.Player2,
		"nextPlayerMove": bc.Game.NextPlayerToMove,
		"gameGrid":       bc.Game.GameGrid,
	})
}
