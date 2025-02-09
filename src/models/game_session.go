package models

import (
	"errors"
	_ "fmt"
	"github.com/google/uuid"
)

var NewUUID = uuid.New

type GameSession struct {
	SessionId        uuid.UUID             `json:"sessionId,omitempty"`
	Player1          string                `json:"player1,omitempty"`
	Player2          string                `json:"player2,omitempty"`
	GameGrid         [3][3]TicTacToeSymbol `json:"gameGrid,omitempty"`
	NextPlayerToMove string                `json:"nextPlayerMove,omitempty"`
}

func NewGameSession(player1 string) *GameSession {
	var grid [3][3]TicTacToeSymbol
	sessionId := NewUUID()

	return &GameSession{
		SessionId:        sessionId,
		Player1:          player1,
		Player2:          "",
		GameGrid:         grid,
		NextPlayerToMove: player1,
	}
}

func (session *GameSession) AddPlayerTwo(player2 string) error {
	if session.Player2 != "" {
		return errors.New("player2 is already set in the session")
	}

	session.Player2 = player2
	return nil
}
