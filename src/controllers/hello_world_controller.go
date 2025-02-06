package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BoardController struct {
	Board [3][3]string
}

func NewBoardController(symbol string) *BoardController {
	controller := &BoardController{}
	controller.Board[2][2] = symbol
	return controller
}

func (bc *BoardController) HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": bc.Board,
	})
}
