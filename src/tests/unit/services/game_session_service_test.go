package services

import (
	"github.com/google/uuid"
	"regexp"
	"src/src/database"
	"src/src/models"
	_ "src/src/models"
	"src/src/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

const mockUuidValue string = "123e4567-e89b-12d3-a456-426614174000"
const playerOne string = "player1"

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
}

func TestRetrieveTicTacToeGameSession(t *testing.T) {
	createGameSessionInDatabase()
	session, _ := services.RetrieveTicTacToeGameSession(mockUuidValue)

	assert.NotNil(t, session)
	assert.Equal(t, playerOne, session.Player1)

	clearGameSessionDatabase()
}

func TestCreateTicTacToeGameSession_ValidUUID(t *testing.T) {
	session := services.CreateTicTacToeGameSession(playerOne)
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	assert.Regexp(t, uuidRegex, session.SessionId)
}

func createGameSessionInDatabase() {
	playerName := playerOne
	session := models.NewGameSession(playerName)
	session.SessionId = mockUuidValue
	db := database.GetInstance()
	db.StoreSession(*session)
}

func clearGameSessionDatabase() {
	db := database.GetInstance()
	db.Clear()
}
