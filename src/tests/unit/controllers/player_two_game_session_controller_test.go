package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"src/src/controllers"
)

func setupPlayerTwoGameSessionRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	hc := controllers.NewPlayerTwoGameSessionController()
	r.POST("/game-session/:gameSessionId/players", hc.PlayerTwoGameSessionControllerHandler)

	return r
}

func TestAddPlayerTwoToGameSessionHandler(t *testing.T) {
	CreateGameSessionInDatabase()

	r := setupPlayerTwoGameSessionRouter()

	jsonBody := `{"player2": "John"}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/game-session/00000000-0000-0000-0000-000000000000/players", bodyReader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t,
		`{
		"sessionId": "00000000-0000-0000-0000-000000000000",
        "player1": "Alice",
        "player2": "John",
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

func TestAddPlayerTwoToGameSessionHandler_whenSessionDoesNotExist(t *testing.T) {
	r := setupPlayerTwoGameSessionRouter()

	jsonBody := `{"player2": "John"}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/game-session/00000000-0000-0000-0000-000000000000/players", bodyReader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t,
		`{"error":"session not found"}`,
		w.Body.String(),
	)
}

func TestAddPlayerTwoToGameSessionHandler_whenSessionAlreadyHasPlayerTwo(t *testing.T) {
	CreateGameSessionInDatabaseWithPlayerTwo()

	r := setupPlayerTwoGameSessionRouter()

	jsonBody := `{"player2": "John"}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/game-session/00000000-0000-0000-0000-000000000000/players", bodyReader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t,
		`{"error":"player2 is already set in the session"}`,
		w.Body.String(),
	)

	ClearGameSessionDatabase()
}
