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

func setupSubmitMoveRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	hc := controllers.NewSubmitMoveController()
	r.POST("/game-session/:gameSessionId", hc.SubmitMoveControllerHandler)

	return r
}

func TestSubmitMoveControllerHandler_success(t *testing.T) {
	CreateGameSessionInDatabaseWithPlayerTwo()

	r := setupSubmitMoveRouter()

	jsonBody := `{"playerName": "Alice",
					"xAxis": 0, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/game-session/00000000-0000-0000-0000-000000000000", bodyReader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t,
		`{
		"sessionId": "00000000-0000-0000-0000-000000000000",
        "player1": "Alice",
        "player2": "John",
        "nextPlayerMove": "John",
		"gameSessionStatus":"Active",
        "gameGrid": [
            ["Cross", "Empty", "Empty"],
            ["Empty", "Empty", "Empty"],
            ["Empty", "Empty", "Empty"]
        ]
    }`,
		w.Body.String(),
	)

	ClearGameSessionDatabase()
}

func TestSubmitMoveControllerHandler_wrongPlayerName(t *testing.T) {
	CreateGameSessionInDatabaseWithPlayerTwo()

	r := setupSubmitMoveRouter()

	jsonBody := `{"playerName": "John",
					"xAxis": 0, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/game-session/00000000-0000-0000-0000-000000000000", bodyReader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t,
		`{"error":"player submitting the move is not next player to move"}`,
		w.Body.String(),
	)

	ClearGameSessionDatabase()
}
