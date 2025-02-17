package controllers_test

import (
	"github.com/google/uuid"
	"src/src/services"
	"testing"
)

func WithMockedUuid(t *testing.T, mockFunc func() uuid.UUID) {
	services.NewUUID = mockFunc
	t.Cleanup(func() {
		services.NewUUID = uuid.New
	})
}

func MockUUID() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000000")
}
