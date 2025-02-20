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

func setupCreateGameRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	hc := controllers.NewCreateGameController()
	r.POST("/game-session", hc.CreateGameControllerHandler)

	return r
}

func TestCreateGameHandler(t *testing.T) {
	WithMockedUuid(t, MockUUID)

	r := setupCreateGameRouter()

	jsonBody := `{"player1": "John"}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/game-session", bodyReader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t,
		`{
		"sessionId": "00000000-0000-0000-0000-000000000000",
        "player1": "John",
        "nextPlayerMove": "John",
		"gameSessionStatus":"Active",
        "gameGrid": [
            ["Empty", "Empty", "Empty"],
            ["Empty", "Empty", "Empty"],
            ["Empty", "Empty", "Empty"]
        ]
    }`,
		w.Body.String(),
	)
}
