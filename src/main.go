package main

import (
	"net/http"
	"src/src/model"

	"github.com/gin-gonic/gin"
	_ "src/src/model"
)

func main() {
	var symbol = model.Circle.String()
	var board [3][3]string
	board[2][2] = symbol

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": board,
		})
	})

	router.Run(":8080")
}
