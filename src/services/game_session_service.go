package services

import (
	"tic-tac-toe-game/src/database"
	"tic-tac-toe-game/src/models"
)

type GameSessionService struct {
	db *database.InMemoryGameSessionDB
}

func NewGameSessionService(db *database.InMemoryGameSessionDB) *GameSessionService {
	return &GameSessionService{db: db}
}

func (s *GameSessionService) CreateTicTacToeGameSession(playerOne string) *models.GameSession {
	newGameSession := models.NewGameSession(playerOne)
	newGameSession.SessionId = GenerateUuid()
	s.db.StoreSession(*newGameSession)

	return newGameSession
}

func (s *GameSessionService) RetrieveTicTacToeGameSession(gameSessionId string) (models.GameSession, error) {
	gameSession, err := s.db.GetSession(gameSessionId)
	if err != nil {
		return models.GameSession{}, err
	}
	return gameSession, nil
}

func (s *GameSessionService) AddPlayerTwoToGameSession(gameSessionId string, playerTwo string) (models.GameSession, error) {
	gameSession, err := s.db.GetSession(gameSessionId)
	if err != nil {
		return models.GameSession{}, err
	}
	err = gameSession.AddPlayerTwo(playerTwo)
	if err != nil {
		return models.GameSession{}, err
	}
	updatedSession, err := s.db.UpdateSession(gameSessionId, gameSession)
	if err != nil {
		return models.GameSession{}, err
	}
	return updatedSession, nil
}

func (s *GameSessionService) AddMoveToGameSession(gameSessionId string, playerName string, xAxis int, yAxis int) (models.GameSession, error) {
	gameSession, getSessionErr := s.db.GetSession(gameSessionId)
	if getSessionErr != nil {
		return models.GameSession{}, getSessionErr
	}

	updatedSession, err := gameSession.SetSymbolOnBoard(playerName, xAxis, yAxis)

	if err != nil {
		return models.GameSession{}, err
	}

	updatedSavedSession, err := s.db.UpdateSession(gameSessionId, *updatedSession)
	if err != nil {
		return models.GameSession{}, err
	}

	return updatedSavedSession, nil
}
