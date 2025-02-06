package model

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
