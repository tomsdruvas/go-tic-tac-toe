package models_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	. "tic-tac-toe-game/src/models"
)

func TestGameSessionStatusString(t *testing.T) {
	tests := []struct {
		status   GameSessionStatus
		expected string
	}{
		{Active, "Active"},
		{Finished, "Finished"},
		{Draw, "Draw"},
		{GameSessionStatus(999), "Unknown"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			assert.Equal(t, test.expected, test.status.String())
		})
	}
}
