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

func TestAddingPlayerTwoToGameSession_whenPlayerTwoAlreadyExists(t *testing.T) {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	_ = session.AddPlayerTwo(playerTwoName)
	err := session.AddPlayerTwo(playerTwoName)
	assert.NotNil(t, err)
	assert.Equal(t, "player2 is already set in the session", err.Error())
}

func TestSetSymbolOnGameSession_whenSubmittingPlayerIsNotNextToPlay(t *testing.T) {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	session.NextPlayerToMove = playerOneName
	session.Player2 = playerTwoName
	_, err := session.SetSymbolOnBoard(playerTwoName, 1, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "player submitting the move is not next player to move", err.Error())
}

func TestSetSymbolOnGameSession_whenSessionIsNotInActiveState(t *testing.T) {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	session.NextPlayerToMove = playerOneName
	session.Player2 = playerTwoName
	session.GameSessionStatus = models.Finished
	_, err := session.SetSymbolOnBoard(playerOneName, 1, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "game session not in Active state", err.Error())
}

func TestSetSymbolOnGameSession_whenSlotAlreadyHasSymbol(t *testing.T) {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	session.NextPlayerToMove = playerOneName
	session.Player2 = playerTwoName
	session.GameGrid[1][1] = models.Cross
	_, err := session.SetSymbolOnBoard(playerOneName, 1, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "the selected slot is already occupied, choose a different slot", err.Error())
}

func TestSetSymbolOnGameSession_whenInputPlayerNameNotPresentInSession(t *testing.T) {
	playerOneName := "Alice"
	session := models.NewGameSession(playerOneName)
	_, err := session.SetSymbolOnBoard("randomName", 1, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "player two not joined", err.Error())
}

func TestSetSymbolOnGameSession_shouldSwitchToPlayerTwoAsNextToPlay(t *testing.T) {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	session.NextPlayerToMove = playerOneName
	session.Player2 = playerTwoName
	_, err := session.SetSymbolOnBoard(playerOneName, 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, playerTwoName, session.NextPlayerToMove)
	assert.Equal(t, models.Active, session.GameSessionStatus)
}

func TestSetSymbolOnGameSession_shouldSwitchToPlayerOneAsNextToPlay(t *testing.T) {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	session.NextPlayerToMove = playerTwoName
	session.Player2 = playerTwoName
	_, err := session.SetSymbolOnBoard(playerTwoName, 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, playerOneName, session.NextPlayerToMove)
}

func TestSetSymbolOnGameSession_whenGameEndUpInFinishedState(t *testing.T) {
	tests := []struct {
		name             string
		inputGameSession *models.GameSession
		nextMove         [2]int
	}{
		{
			name:             "Column has 3 consecutive",
			inputGameSession: createGameSessionWithSymbolsOnGrid([2]int{0, 0}, [2]int{0, 1}),
			nextMove:         [2]int{0, 2},
		},
		{
			name:             "Row has 3 consecutive",
			inputGameSession: createGameSessionWithSymbolsOnGrid([2]int{0, 0}, [2]int{1, 0}),
			nextMove:         [2]int{2, 0},
		},
		{
			name:             "Diagonal has 3 consecutive",
			inputGameSession: createGameSessionWithSymbolsOnGrid([2]int{0, 0}, [2]int{1, 1}),
			nextMove:         [2]int{2, 2},
		},
		{
			name:             "Anti-Diagonal has 3 consecutive",
			inputGameSession: createGameSessionWithSymbolsOnGrid([2]int{0, 2}, [2]int{1, 1}),
			nextMove:         [2]int{2, 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			returnedGameSession, err := tt.inputGameSession.SetSymbolOnBoard("Alice", tt.nextMove[0], tt.nextMove[1])
			assert.NoError(t, err)
			assert.Equal(t, models.Finished, returnedGameSession.GameSessionStatus)
			assert.Equal(t, "Alice", returnedGameSession.NextPlayerToMove)
			assert.Equal(t, "Alice", returnedGameSession.Winner)
		})
	}
}

func TestSetSymbolOnGameSession_whenNextMoveEndsInDraw(t *testing.T) {
	session := createGameSessionGridWithEightSlotsFilled()
	_, err := session.SetSymbolOnBoard("Alice", 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, models.Draw, session.GameSessionStatus)
}

func createGameSessionWithSymbolsOnGrid(first [2]int, second [2]int) *models.GameSession {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	_ = session.AddPlayerTwo(playerTwoName)
	session.NextPlayerToMove = playerOneName
	session.GameGrid[first[0]][first[1]] = models.Cross
	session.GameGrid[second[0]][second[1]] = models.Cross
	return session
}

func createGameSessionGridWithEightSlotsFilled() *models.GameSession {
	playerOneName := "Alice"
	playerTwoName := "John"
	session := models.NewGameSession(playerOneName)
	_ = session.AddPlayerTwo(playerTwoName)
	session.NextPlayerToMove = playerOneName

	session.GameGrid = [3][3]models.TicTacToeSymbol{
		{models.Cross, models.Circle, models.Empty},
		{models.Circle, models.Cross, models.Circle},
		{models.Circle, models.Cross, models.Circle},
	}

	return session
}
