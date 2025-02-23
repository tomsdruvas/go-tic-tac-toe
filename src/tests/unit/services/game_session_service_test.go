package services

import (
	"github.com/google/uuid"
	"regexp"
	"testing"
	"tic-tac-toe-game/src/database"
	"tic-tac-toe-game/src/models"
	_ "tic-tac-toe-game/src/models"
	"tic-tac-toe-game/src/services"

	"github.com/stretchr/testify/assert"
)

const mockUuidValue string = "123e4567-e89b-12d3-a456-426614174000"
const playerOne string = "player1"
const playerTwo string = "player2"

func mockUUID() uuid.UUID {
	mockedUUID, _ := uuid.Parse(mockUuidValue)
	return mockedUUID
}

func TestCreateTicTacToeGameSession(t *testing.T) {
	originalUUID := services.NewUUID
	services.NewUUID = mockUUID
	defer func() { services.NewUUID = originalUUID }()

	session := services.CreateTicTacToeGameSession(playerOne)

	assert.NotNil(t, session)
	expectedUUID := mockUuidValue
	assert.Equal(t, expectedUUID, session.SessionId)
	assert.Equal(t, playerOne, session.Player1)
	assert.Equal(t, models.Active, session.GameSessionStatus)
}

func TestRetrieveTicTacToeGameSession(t *testing.T) {
	createGameSessionInDatabase()
	session, _ := services.RetrieveTicTacToeGameSession(mockUuidValue)

	assert.NotNil(t, session)
	assert.Equal(t, playerOne, session.Player1)
	assert.Equal(t, models.Active, session.GameSessionStatus)

	clearGameSessionDatabase()
}

func TestCreateTicTacToeGameSession_ValidUUID(t *testing.T) {
	session := services.CreateTicTacToeGameSession(playerOne)
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	assert.Regexp(t, uuidRegex, session.SessionId)
	clearGameSessionDatabase()
}

func TestAddPlayerTwoToGameSession(t *testing.T) {
	createGameSessionInDatabase()
	updatedGameSession, err := services.AddPlayerTwoToGameSession(mockUuidValue, playerTwo)
	assert.Nil(t, err)
	assert.NotNil(t, updatedGameSession)
	assert.Equal(t, playerOne, updatedGameSession.Player1)
	assert.Equal(t, playerTwo, updatedGameSession.Player2)
	assert.Equal(t, models.Active, updatedGameSession.GameSessionStatus)

	retrievedSession, err := services.RetrieveTicTacToeGameSession(mockUuidValue)
	assert.Nil(t, err)
	assert.NotNil(t, retrievedSession)
	assert.Equal(t, playerOne, retrievedSession.Player1)
	assert.Equal(t, playerTwo, retrievedSession.Player2)
	assert.Equal(t, models.Active, retrievedSession.GameSessionStatus)
	clearGameSessionDatabase()
}

func TestAddPlayerTwoToGameSession_sessionNotFound(t *testing.T) {
	_, err := services.AddPlayerTwoToGameSession(mockUuidValue, playerTwo)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "session not found")
}

func TestAddPlayerTwoToGameSession_playerTwoAlreadyExists(t *testing.T) {
	createGameSessionInDatabaseWith2Players()
	_, err := services.AddPlayerTwoToGameSession(mockUuidValue, playerTwo)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "player2 is already set in the session")
	clearGameSessionDatabase()
}

func TestAddMoveOnGameSession(t *testing.T) {
	createGameSessionInDatabaseWith2Players()
	gameSession, err := services.AddMoveToGameSession(mockUuidValue, playerOne, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, playerTwo, gameSession.NextPlayerToMove)
	clearGameSessionDatabase()
}

func TestAddMoveOnGameSession_withPlayerTwoWhenPlayerOneShouldMakeMove(t *testing.T) {
	createGameSessionInDatabaseWith2Players()
	_, err := services.AddMoveToGameSession(mockUuidValue, playerTwo, 0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "player submitting the move is not next player to move")
	clearGameSessionDatabase()
}

func TestAddMoveOnGameSession_whenSessionNotFound(t *testing.T) {
	_, err := services.AddMoveToGameSession(mockUuidValue, playerTwo, 0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "session not found")
}

func createGameSessionInDatabase() {
	playerName := playerOne
	session := models.NewGameSession(playerName)
	session.SessionId = mockUuidValue
	db := database.GetInstance()
	db.StoreSession(*session)
}

func createGameSessionInDatabaseWith2Players() {
	playerName := playerOne
	session := models.NewGameSession(playerName)
	session.Player2 = playerTwo
	session.SessionId = mockUuidValue
	db := database.GetInstance()
	db.StoreSession(*session)
}

func clearGameSessionDatabase() {
	db := database.GetInstance()
	db.Clear()
}
