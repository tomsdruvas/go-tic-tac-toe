package services

import (
	"github.com/google/uuid"
	"regexp"
	_ "src/src/models"
	"src/src/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

const MockUuidValue string = "123e4567-e89b-12d3-a456-426614174000"
const PlayerOne string = "player1"

func mockUUID() uuid.UUID {
	mockedUUID, _ := uuid.Parse(MockUuidValue)
	return mockedUUID
}

func TestCreateTicTacToeGameSession(t *testing.T) {
	originalUUID := services.NewUUID
	services.NewUUID = mockUUID
	defer func() { services.NewUUID = originalUUID }()

	session := services.CreateTicTacToeGameSession(PlayerOne)

	assert.NotNil(t, session)
	expectedUUID := MockUuidValue
	assert.Equal(t, expectedUUID, session.SessionId)
	assert.Equal(t, PlayerOne, session.Player1)
}

func TestCreateTicTacToeGameSession_ValidUUID(t *testing.T) {
	session := services.CreateTicTacToeGameSession(PlayerOne)
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	assert.Regexp(t, uuidRegex, session.SessionId)
}
