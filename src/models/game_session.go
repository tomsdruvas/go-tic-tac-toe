package models

import _ "fmt"

type GameSession struct {
	Player1          string
	Player2          string
	GameGrid         [3][3]string
	NextPlayerToMove string
}

func NewGameSession(player1 string) *GameSession {
	var grid [3][3]string

	return &GameSession{
		Player1:          player1,
		Player2:          "",
		GameGrid:         grid,
		NextPlayerToMove: player1,
	}
}
