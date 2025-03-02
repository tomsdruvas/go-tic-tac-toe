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

	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)

	session := gameSessionService.CreateTicTacToeGameSession(playerOne)
	assert.NotNil(t, session)
	expectedUUID := mockUuidValue
	assert.Equal(t, expectedUUID, session.SessionId)
	assert.Equal(t, playerOne, session.Player1)
	assert.Equal(t, models.Active, session.GameSessionStatus)
}

func TestRetrieveTicTacToeGameSession(t *testing.T) {
	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)

	createGameSessionInDatabase(db)
	session, _ := gameSessionService.RetrieveTicTacToeGameSession(mockUuidValue)

	assert.NotNil(t, session)
	assert.Equal(t, playerOne, session.Player1)
	assert.Equal(t, models.Active, session.GameSessionStatus)
}

func TestCreateTicTacToeGameSession_ValidUUID(t *testing.T) {
	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)

	session := gameSessionService.CreateTicTacToeGameSession(playerOne)
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	assert.Regexp(t, uuidRegex, session.SessionId)

}

func TestAddPlayerTwoToGameSession(t *testing.T) {
	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)
	createGameSessionInDatabase(db)
	updatedGameSession, err := gameSessionService.AddPlayerTwoToGameSession(mockUuidValue, playerTwo)
	assert.Nil(t, err)
	assert.NotNil(t, updatedGameSession)
	assert.Equal(t, playerOne, updatedGameSession.Player1)
	assert.Equal(t, playerTwo, updatedGameSession.Player2)
	assert.Equal(t, models.Active, updatedGameSession.GameSessionStatus)

	retrievedSession, err := gameSessionService.RetrieveTicTacToeGameSession(mockUuidValue)
	assert.Nil(t, err)
	assert.NotNil(t, retrievedSession)
	assert.Equal(t, playerOne, retrievedSession.Player1)
	assert.Equal(t, playerTwo, retrievedSession.Player2)
	assert.Equal(t, models.Active, retrievedSession.GameSessionStatus)
}

func TestAddPlayerTwoToGameSession_sessionNotFound(t *testing.T) {
	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)
	_, err := gameSessionService.AddPlayerTwoToGameSession(mockUuidValue, playerTwo)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "session not found")
}

func TestAddPlayerTwoToGameSession_playerTwoAlreadyExists(t *testing.T) {
	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)
	createGameSessionInDatabaseWith2Players(db)
	_, err := gameSessionService.AddPlayerTwoToGameSession(mockUuidValue, playerTwo)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "player2 is already set in the session")
}

func TestAddMoveOnGameSession(t *testing.T) {
	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)
	createGameSessionInDatabaseWith2Players(db)
	gameSession, err := gameSessionService.AddMoveToGameSession(mockUuidValue, playerOne, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, playerTwo, gameSession.NextPlayerToMove)
}

func TestAddMoveOnGameSession_withPlayerTwoWhenPlayerOneShouldMakeMove(t *testing.T) {
	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)
	createGameSessionInDatabaseWith2Players(db)
	_, err := gameSessionService.AddMoveToGameSession(mockUuidValue, playerTwo, 0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "player submitting the move is not next player to move")
}

func TestAddMoveOnGameSession_whenSessionNotFound(t *testing.T) {
	db := database.NewInMemoryGameSessionDB()
	gameSessionService := services.NewGameSessionService(db)
	_, err := gameSessionService.AddMoveToGameSession(mockUuidValue, playerTwo, 0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "session not found")
}

func createGameSessionInDatabase(db *database.InMemoryGameSessionDB) {
	playerName := playerOne
	session := models.NewGameSession(playerName)
	session.SessionId = mockUuidValue
	db.StoreSession(*session)
}

func createGameSessionInDatabaseWith2Players(db *database.InMemoryGameSessionDB) {
	playerName := playerOne
	session := models.NewGameSession(playerName)
	session.Player2 = playerTwo
	session.SessionId = mockUuidValue
	db.StoreSession(*session)
}
