package models_test

import (
	"github.com/stretchr/testify/assert"
	. "src/src/models"
	"testing"
)

func TestTicTacToeSymbolString(t *testing.T) {
	tests := []struct {
		symbol   TicTacToeSymbol
		expected string
	}{
		{Empty, "Empty"},
		{Circle, "Circle"},
		{Cross, "Cross"},
		{TicTacToeSymbol(999), "Unknown"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			assert.Equal(t, test.expected, test.symbol.String(), "String representation does not match")
		})
	}
}
