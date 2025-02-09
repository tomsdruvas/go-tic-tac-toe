package controllers_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestHealthCheckIntegration(t *testing.T) {
	resp, err := http.Get(TestServerURL + "/health")
	assert.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bytes.TrimSpace(bodyBytes))
	expectedBody := `{"message":"Service is up and running","status":"healthy"}`

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, expectedBody, bodyString, "Unexpected response body")
}
