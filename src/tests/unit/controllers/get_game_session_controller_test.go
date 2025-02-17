package controllers_test

import (
	_ "github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"src/src/database"
	"src/src/models"
	_ "src/src/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"src/src/controllers"
)

func setupGetGameSessionRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	hc := controllers.NewGetGameSessionController()
	r.GET("/game-session/:gameSessionId", hc.GetGameSessionControllerHandler)

	return r
}

func TestGetGameSessionSuccessfulHandler(t *testing.T) {
	createGameSessionInDatabase()
	r := setupGetGameSessionRouter()

	req, _ := http.NewRequest("GET", "/game-session/00000000-0000-0000-0000-000000000000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t,
		`{
		"sessionId": "00000000-0000-0000-0000-000000000000",
        "player1": "Alice",
        "nextPlayerMove": "Alice",
        "gameGrid": [
            [0, 0, 0],
            [0, 0, 0],
            [0, 0, 0]
        ]
    }`,
		w.Body.String(),
	)

	clearGameSessionDatabase()
}

func TestGetGameSessionNotFoundHandler(t *testing.T) {
	r := setupGetGameSessionRouter()

	req, _ := http.NewRequest("GET", "/game-session/00000000-0000-0000-0000-000000000000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t,
		`{"error": "Game session not found"}`,
		w.Body.String(),
	)
}

func createGameSessionInDatabase() {
	playerName := "Alice"
	session := models.NewGameSession(playerName)
	session.SessionId = "00000000-0000-0000-0000-000000000000"
	db := database.GetInstance()
	db.StoreSession(*session)
}

func clearGameSessionDatabase() {
	db := database.GetInstance()
	db.Clear()
}
