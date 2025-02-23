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

func (session *GameSession) SetSymbolOnBoard(playerName string, xAxis, yAxis int) (*GameSession, error) {
	if !session.hasPlayerTwoJoined() {
		return nil, errors.New("player two not joined")
	}

	if !session.isNextPlayer(playerName) {
		return nil, errors.New("player submitting the move is not next player to move")
	}

	if !session.isActive() {
		return nil, errors.New("game session not in Active state")
	}

	symbol := session.getSymbolForPlayer(playerName)

	if !session.isSlotEmpty(xAxis, yAxis) {
		return nil, errors.New("the selected slot is already occupied, choose a different slot")
	}

	session.GameGrid[xAxis][yAxis] = symbol

	session.calculateGameSessionStatus()

	if session.GameSessionStatus == Active {
		if session.Player1 == playerName {
			session.NextPlayerToMove = session.Player2
		} else {
			session.NextPlayerToMove = session.Player1
		}
	}
	return session, nil
}

func (session *GameSession) hasPlayerTwoJoined() bool {
	return session.Player2 != ""
}

func (session *GameSession) isNextPlayer(playerName string) bool {
	return playerName == session.NextPlayerToMove
}

func (session *GameSession) isActive() bool {
	return session.GameSessionStatus == Active
}

func (session *GameSession) getSymbolForPlayer(playerName string) TicTacToeSymbol {
	if playerName == session.Player1 {
		return Cross
	}
	return Circle
}

func (session *GameSession) isSlotEmpty(xAxis, yAxis int) bool {
	return session.GameGrid[xAxis][yAxis] == Empty
}

func (session *GameSession) calculateGameSessionStatus() {
	if session.hasWinningRow() || session.hasWinningColumn() || session.hasWinningDiagonal() {
		session.GameSessionStatus = Finished
		session.Winner = session.NextPlayerToMove
		return
	}

	if session.isDraw() {
		session.GameSessionStatus = Draw
	}
}

func (session *GameSession) hasWinningRow() bool {
	for i := 0; i < 3; i++ {
		if session.GameGrid[i][0] != Empty &&
			session.GameGrid[i][0] == session.GameGrid[i][1] &&
			session.GameGrid[i][1] == session.GameGrid[i][2] {
			return true
		}
	}
	return false
}

func (session *GameSession) hasWinningColumn() bool {
	for i := 0; i < 3; i++ {
		if session.GameGrid[0][i] != Empty &&
			session.GameGrid[0][i] == session.GameGrid[1][i] &&
			session.GameGrid[1][i] == session.GameGrid[2][i] {
			return true
		}
	}
	return false
}

func (session *GameSession) hasWinningDiagonal() bool {
	grid := session.GameGrid
	//diagonal
	if grid[0][0] != Empty && grid[0][0] == grid[1][1] && grid[1][1] == grid[2][2] {
		return true
	}
	//anti-diag
	if grid[0][2] != Empty && grid[0][2] == grid[1][1] && grid[1][1] == grid[2][0] {
		return true
	}
	return false
}

func (session *GameSession) isDraw() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if session.GameGrid[i][j] == Empty {
				return false
			}
		}
	}
	return true
}
