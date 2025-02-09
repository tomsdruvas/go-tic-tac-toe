package controllers_test

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"src/src/models"
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
	r := setupCreateGameRouter()
	mockUuidGeneration(t)

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

func mockUuidGeneration(t *testing.T) {
	origNewUUID := models.NewUUID
	fixedUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	models.NewUUID = func() uuid.UUID {
		return fixedUUID
	}

	t.Cleanup(func() {
		models.NewUUID = origNewUUID
	})
}
