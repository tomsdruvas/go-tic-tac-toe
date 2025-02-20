package models

import "encoding/json"

type GameSessionStatus int

const (
	Active GameSessionStatus = iota
	Finished
	Draw
)

func (s GameSessionStatus) String() string {
	switch s {
	case Active:
		return "Active"
	case Finished:
		return "Finished"
	case Draw:
		return "Draw"
	default:
		return "Unknown"
	}
}

func (s GameSessionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
