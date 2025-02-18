package controllers_test

import (
	"github.com/google/uuid"
	"src/src/database"
	"src/src/models"
	"src/src/services"
	"testing"
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

func CreateGameSessionInDatabase() {
	playerName := "Alice"
	session := models.NewGameSession(playerName)
	session.SessionId = "00000000-0000-0000-0000-000000000000"
	db := database.GetInstance()
	db.StoreSession(*session)
}

func CreateGameSessionInDatabaseWithPlayerTwo() {
	playerName := "Alice"
	session := models.NewGameSession(playerName)
	session.Player2 = "John"
	session.SessionId = "00000000-0000-0000-0000-000000000000"
	db := database.GetInstance()
	db.StoreSession(*session)
}

func ClearGameSessionDatabase() {
	db := database.GetInstance()
	db.Clear()
}
