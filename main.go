package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"src/src/controllers"
	"src/src/models"
)

func main() {
	symbol := models.Circle.String()

	boardController := controllers.NewBoardController(symbol)

	router := gin.Default()

	router.GET("/hello", boardController.HelloHandler)

	router.Run(":8080")
}
