package services

import "github.com/google/uuid"

var NewUUID = uuid.New

func GenerateUuid() string {
	sessionId := NewUUID()
	return sessionId.String()
}
