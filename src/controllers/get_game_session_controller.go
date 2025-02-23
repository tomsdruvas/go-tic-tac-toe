package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"src/src/services"
)

type GetGameSessionController struct {
}

func NewGetGameSessionController() *GetGameSessionController {
	return &GetGameSessionController{}
}

func (bc *GetGameSessionController) GetGameSessionControllerHandler(c *gin.Context) {
	gameSessionId := c.Param("gameSessionId")
	gameSession, err := services.RetrieveTicTacToeGameSession(gameSessionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game session not found"})
		return
	}

	c.JSON(http.StatusOK, gameSession)
}
