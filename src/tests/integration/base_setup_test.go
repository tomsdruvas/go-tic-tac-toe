package controllers_test

import (
	"os"
	. "src/src/tests/integration/setup"
	"testing"
)

const TestServerURL = "http://localhost:8081"

func TestMain(m *testing.M) {
	SetupDocker()
	code := m.Run()
	StopDockerCompose()
	os.Exit(code)
}
