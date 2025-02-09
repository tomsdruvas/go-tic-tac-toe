package models

import (
	"errors"
	_ "fmt"
)

type GameSession struct {
	Player1          string
	Player2          string
	GameGrid         [3][3]TicTacToeSymbol
	NextPlayerToMove string
}

func NewGameSession(player1 string) *GameSession {
	var grid [3][3]TicTacToeSymbol

	return &GameSession{
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
