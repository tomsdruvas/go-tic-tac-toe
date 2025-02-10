package services

import (
	"regexp"
	"src/src/services"
	"testing"
)

var uuidRegex = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

func TestGenerateUuid(t *testing.T) {
	uuid := services.GenerateUuid()

	if !uuidRegex.MatchString(uuid) {
		t.Errorf("GenerateUuid() returned an invalid UUID: %s", uuid)
	}
}
