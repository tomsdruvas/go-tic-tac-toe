package controllers_test

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"src/src/services"
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
	r.POST("/create-game", hc.CreateGameControllerHandler)

	return r
}

func TestCreateGameHandler(t *testing.T) {
	withMockedUUID(t, mockUUID)

	r := setupCreateGameRouter()

	jsonBody := `{"player1": "John"}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/create-game", bodyReader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t,
		`{
		"sessionId": "00000000-0000-0000-0000-000000000000",
        "player1": "John",
        "nextPlayerMove": "John",
        "gameGrid": [
            [0, 0, 0],
            [0, 0, 0],
            [0, 0, 0]
        ]
    }`,
		w.Body.String(),
	)
}

func mockUUID() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000000")
}

func withMockedUUID(t *testing.T, mockFunc func() uuid.UUID) {
	services.NewUUID = mockFunc
	t.Cleanup(func() {
		services.NewUUID = uuid.New
	})
}
