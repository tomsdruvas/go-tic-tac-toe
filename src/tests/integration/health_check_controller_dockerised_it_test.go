package controllers_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	. "src/src/tests/integration/setup"
	"testing"
)

const testServerURL = "http://localhost:8081"

func TestMain(m *testing.M) {
	SetupDocker()
	code := m.Run()
	StopDockerCompose()
	os.Exit(code)
}

func TestHealthCheckIntegration(t *testing.T) {
	resp, err := http.Get(testServerURL + "/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
