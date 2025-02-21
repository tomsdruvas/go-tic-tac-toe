package models

import (
	"errors"
	_ "fmt"
)

type GameSession struct {
	SessionId         string                `json:"sessionId,omitempty"`
	Player1           string                `json:"player1,omitempty"`
	Player2           string                `json:"player2,omitempty"`
	GameGrid          [3][3]TicTacToeSymbol `json:"gameGrid,omitempty"`
	NextPlayerToMove  string                `json:"nextPlayerMove,omitempty"`
	GameSessionStatus GameSessionStatus     `json:"gameSessionStatus"`
	Winner            string                `json:"winner,omitempty"`
}

func NewGameSession(player1 string) *GameSession {
	var grid [3][3]TicTacToeSymbol

	return &GameSession{
		Player1:           player1,
		GameGrid:          grid,
		NextPlayerToMove:  player1,
		GameSessionStatus: Active,
	}
}

func (session *GameSession) AddPlayerTwo(player2 string) error {
	if session.Player2 != "" {
		return errors.New("player2 is already set in the session")
	}

	session.Player2 = player2
	return nil
}
