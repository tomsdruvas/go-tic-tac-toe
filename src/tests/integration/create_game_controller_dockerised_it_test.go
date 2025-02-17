package integration_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	_ "src/src/tests/integration/setup"
	"strings"
	"testing"
)

var currentSessionId string

func TestTicTacToeFullFlow(t *testing.T) {
	createGameSessionIntegration(t)
	retrieveGameSessionIntegration(t)
}

func TestRetrieveNotFoundGameSessionIntegration(t *testing.T) {
	resp, err := http.Get(TestServerURL + "/game-session/1234")
	assert.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func createGameSessionIntegration(t *testing.T) {
	jsonBody := `{"player1": "John"}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session", "application/json", bodyReader)
	assert.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	sessionId, ok := response["sessionId"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, sessionId)
	currentSessionId = sessionId

	assert.Equal(t, "John", response["player1"])
	assert.Equal(t, "John", response["nextPlayerMove"])

	expectedGrid := []interface{}{
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
	}
	assert.Equal(t, expectedGrid, response["gameGrid"])
}

func retrieveGameSessionIntegration(t *testing.T) {
	resp, err := http.Get(TestServerURL + "/game-session/" + currentSessionId)
	assert.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	sessionId, ok := response["sessionId"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, sessionId)
	currentSessionId = sessionId

	assert.Equal(t, "John", response["player1"])
	assert.Equal(t, "John", response["nextPlayerMove"])

	expectedGrid := []interface{}{
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
	}
	assert.Equal(t, expectedGrid, response["gameGrid"])
}
