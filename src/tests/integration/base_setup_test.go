package integration_test

import (
	"os"
	"src/src/tests/integration/setup"
	"testing"
)

const TestServerURL = "http://localhost:8081"

func TestMain(m *testing.M) {
	setup.SpinUpDocker()
	code := m.Run()
	setup.StopDockerCompose()
	os.Exit(code)
}
