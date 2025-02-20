package controllers_test

import (
	_ "github.com/google/uuid"
	"net/http"
	"net/http/httptest"
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
	CreateGameSessionInDatabase()
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
		"gameSessionStatus":"Active",
        "gameGrid": [
            ["Empty", "Empty", "Empty"],
            ["Empty", "Empty", "Empty"],
            ["Empty", "Empty", "Empty"]
        ]
    }`,
		w.Body.String(),
	)

	ClearGameSessionDatabase()
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
