package services

import (
	"src/src/database"
	"src/src/models"
)

func CreateTicTacToeGameSession(playerOne string) *models.GameSession {
	newGameSession := models.NewGameSession(playerOne)
	newGameSession.SessionId = GenerateUuid()
	databaseInstance := database.GetInstance()
	databaseInstance.StoreSession(*newGameSession)

	return newGameSession
}

func RetrieveTicTacToeGameSession(gameSessionId string) (models.GameSession, error) {
	databaseInstance := database.GetInstance()
	gameSession, err := databaseInstance.GetSession(gameSessionId)
	if err != nil {
		return models.GameSession{}, err
	}
	return gameSession, nil
}
