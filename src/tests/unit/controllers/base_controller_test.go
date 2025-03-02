package controllers_test

import (
	"github.com/google/uuid"
	"testing"
	"tic-tac-toe-game/src/database"
	"tic-tac-toe-game/src/models"
	"tic-tac-toe-game/src/services"
)

func WithMockedUuid(t *testing.T, mockFunc func() uuid.UUID) {
	services.NewUUID = mockFunc
	t.Cleanup(func() {
		services.NewUUID = uuid.New
	})
}

func MockUUID() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000000")
}

func CreateGameSessionInDatabase(db *database.InMemoryGameSessionDB) {
	playerName := "Alice"
	session := models.NewGameSession(playerName)
	session.SessionId = "00000000-0000-0000-0000-000000000000"
	db.StoreSession(*session)
}

func CreateGameSessionInDatabaseWithPlayerTwo(db *database.InMemoryGameSessionDB) {
	playerName := "Alice"
	session := models.NewGameSession(playerName)
	session.Player2 = "John"
	session.SessionId = "00000000-0000-0000-0000-000000000000"
	db.StoreSession(*session)
}

func ClearGameSessionDatabase(db *database.InMemoryGameSessionDB) {
	db.Clear()
}
