package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"src/src/models"
	"src/src/services"
)

type SubmitMoveController struct {
	Game *models.GameSession
}

type SubmitMoveRequest struct {
	PlayerName string `json:"playerName" binding:"required"`
	XAxis      *int   `json:"xAxis" binding:"required"`
	YAxis      *int   `json:"yAxis" binding:"required"`
}

func NewSubmitMoveController() *SubmitMoveController {
	return &SubmitMoveController{}
}

func (bc *SubmitMoveController) SubmitMoveControllerHandler(c *gin.Context) {
	var req SubmitMoveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameSessionId := c.Param("gameSessionId")

	session, err := services.AddMoveToGameSession(gameSessionId, req.PlayerName, *req.XAxis, *req.YAxis)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}
