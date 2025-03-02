package controllers_test

import (
	"go.uber.org/dig"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tic-tac-toe-game/src/database"
	"tic-tac-toe-game/src/services"
	"tic-tac-toe-game/src/websockets"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"tic-tac-toe-game/src/controllers"
)

func setupGameSessionRouter(t *testing.T) (*gin.Engine, *database.InMemoryGameSessionDB) {
	gin.SetMode(gin.TestMode)

	container := dig.New()

	assert.NoError(t, container.Provide(database.NewInMemoryGameSessionDB))
	assert.NoError(t, container.Provide(services.NewGameSessionService))
	assert.NoError(t, container.Provide(websockets.NewWebSocketService))
	assert.NoError(t, container.Provide(controllers.NewGameSessionController))

	var controller *controllers.GameSessionController
	var inMemoryDb *database.InMemoryGameSessionDB
	err := container.Invoke(func(c *controllers.GameSessionController, db *database.InMemoryGameSessionDB) {
		controller = c
		inMemoryDb = db
	})
	assert.NoError(t, err)

	r := gin.Default()
	r.POST("/game-session", controller.CreateGameSessionHandler)
	r.GET("/game-session/:gameSessionId", controller.GetGameSessionHandler)
	r.POST("/game-session/:gameSessionId/players", controller.PlayerTwoGameSessionHandler)
	r.POST("/game-session/:gameSessionId/move", controller.SubmitMoveHandler)

	return r, inMemoryDb
}

func TestCreateGameHandler_happyPath(t *testing.T) {
	WithMockedUuid(t, MockUUID)

	r, _ := setupGameSessionRouter(t)

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

func TestGetGameSession_happyPath(t *testing.T) {
	r, db := setupGameSessionRouter(t)
	CreateGameSessionInDatabase(db)

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

	ClearGameSessionDatabase(db)
}

func TestGetGameSessionNotFoundHandler(t *testing.T) {
	r, _ := setupGameSessionRouter(t)

	req, _ := http.NewRequest("GET", "/game-session/00000000-0000-0000-0000-000000000000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t,
		`{"error": "Game session not found"}`,
		w.Body.String(),
	)
}

func TestAddPlayerTwoToGameSessionHandler(t *testing.T) {
	r, db := setupGameSessionRouter(t)
	CreateGameSessionInDatabase(db)

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

	ClearGameSessionDatabase(db)
}

func TestAddPlayerTwoToGameSessionHandler_whenSessionDoesNotExist(t *testing.T) {
	r, _ := setupGameSessionRouter(t)

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
	r, db := setupGameSessionRouter(t)
	CreateGameSessionInDatabaseWithPlayerTwo(db)

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

	ClearGameSessionDatabase(db)
}

func TestSubmitMoveControllerHandler_success(t *testing.T) {
	r, db := setupGameSessionRouter(t)
	CreateGameSessionInDatabaseWithPlayerTwo(db)

	jsonBody := `{"playerName": "Alice",
					"xAxis": 0, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/game-session/00000000-0000-0000-0000-000000000000/move", bodyReader)
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

	ClearGameSessionDatabase(db)
}

func TestSubmitMoveControllerHandler_wrongPlayerName(t *testing.T) {
	r, db := setupGameSessionRouter(t)
	CreateGameSessionInDatabaseWithPlayerTwo(db)

	jsonBody := `{"playerName": "John",
					"xAxis": 0, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/game-session/00000000-0000-0000-0000-000000000000/move", bodyReader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t,
		`{"error":"player submitting the move is not next player to move"}`,
		w.Body.String(),
	)

	ClearGameSessionDatabase(db)
}
