package integration_test

import (
	"os"
	"testing"
	"tic-tac-toe-game/src/tests/integration/setup"
)

const TestServerURL = "http://localhost:8081"

func TestMain(m *testing.M) {
	setup.SpinUpDocker()
	code := m.Run()
	setup.StopDockerCompose()
	os.Exit(code)
}
