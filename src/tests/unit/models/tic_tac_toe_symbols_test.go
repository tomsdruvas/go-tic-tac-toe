package models_test

import (
	"github.com/stretchr/testify/assert"
	"src/src/models"
	"testing"
)

func TestTicTacToeSymbolString(t *testing.T) {
	tests := []struct {
		symbol   models.TicTacToeSymbol
		expected string
	}{
		{models.Empty, "Empty"},
		{models.Circle, "Circle"},
		{models.Cross, "Cross"},
		{models.TicTacToeSymbol(999), "Unknown"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			assert.Equal(t, test.expected, test.symbol.String(), "String representation does not match")
		})
	}
}
