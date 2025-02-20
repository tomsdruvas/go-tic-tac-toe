package models

import "encoding/json"

type TicTacToeSymbol int

const (
	Empty TicTacToeSymbol = iota
	Circle
	Cross
)

func (s TicTacToeSymbol) String() string {
	switch s {
	case Empty:
		return "Empty"
	case Circle:
		return "Circle"
	case Cross:
		return "Cross"
	default:
		return "Unknown"
	}
}

func (s TicTacToeSymbol) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
