package services

import (
	"src/src/models"
)

func CreateTicTacToeGameSession(playerOne string) *models.GameSession {
	newGameSession := models.NewGameSession(playerOne)
	newGameSession.SessionId = GenerateUuid()

	return newGameSession
}
