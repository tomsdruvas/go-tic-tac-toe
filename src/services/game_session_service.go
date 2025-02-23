package services

import (
	"tic-tac-toe-game/src/database"
	"tic-tac-toe-game/src/models"
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

func AddPlayerTwoToGameSession(gameSessionId string, playerTwo string) (models.GameSession, error) {
	databaseInstance := database.GetInstance()
	gameSession, err := databaseInstance.GetSession(gameSessionId)
	if err != nil {
		return models.GameSession{}, err
	}
	err = gameSession.AddPlayerTwo(playerTwo)
	if err != nil {
		return models.GameSession{}, err
	}
	updatedSession, err := databaseInstance.UpdateSession(gameSessionId, gameSession)
	if err != nil {
		return models.GameSession{}, err
	}
	return updatedSession, nil
}

func AddMoveToGameSession(gameSessionId string, playerName string, xAxis int, yAxis int) (models.GameSession, error) {
	databaseInstance := database.GetInstance()
	gameSession, getSessionErr := databaseInstance.GetSession(gameSessionId)
	if getSessionErr != nil {
		return models.GameSession{}, getSessionErr
	}

	updatedSession, err := gameSession.SetSymbolOnBoard(playerName, xAxis, yAxis)

	if err != nil {
		return models.GameSession{}, err
	}

	updatedSavedSession, err := databaseInstance.UpdateSession(gameSessionId, *updatedSession)
	if err != nil {
		return models.GameSession{}, err
	}

	return updatedSavedSession, nil
}
