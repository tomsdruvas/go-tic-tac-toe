package models

import (
	"github.com/stretchr/testify/assert"
	"src/src/models"
	"testing"
)

func TestNewGameSession(t *testing.T) {
	playerName := "Alice"
	session := models.NewGameSession(playerName)

	assert.NotNil(t, session)
	assert.Equal(t, playerName, session.Player1)
	assert.Equal(t, "", session.Player2)
	assert.Equal(t, playerName, session.NextPlayerToMove)
	expectedGrid := [3][3]models.TicTacToeSymbol{}
	assert.Equal(t, expectedGrid, session.GameGrid)
}

func TestAddingPlayerTwoToGameSession(t *testing.T) {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	_ = session.AddPlayerTwo(playerTwoName)
	assert.Equal(t, playerTwoName, session.Player2)
}
