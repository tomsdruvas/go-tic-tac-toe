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
	assertBodyWithOnePlayer(retrieveGameSessionWithResponseBodyIntegration(t))
	addPlayerTwoToGameSessionIntegration(t)
	assertBodyWithTwoPlayer(retrieveGameSessionWithResponseBodyIntegration(t))
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

func retrieveGameSessionWithResponseBodyIntegration(t *testing.T) (tt *testing.T, resp *http.Response, err error) {
	resp, err = http.Get(TestServerURL + "/game-session/" + currentSessionId)
	return t, resp, err
}

func assertBodyWithOnePlayer(t *testing.T, resp *http.Response, err error) {
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
	assert.NotContains(t, response, "player2")

	expectedGrid := []interface{}{
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
	}
	assert.Equal(t, expectedGrid, response["gameGrid"])
}

func assertBodyWithTwoPlayer(t *testing.T, resp *http.Response, err error) {
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
	assert.Equal(t, "Alice", response["player2"])

	expectedGrid := []interface{}{
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
	}
	assert.Equal(t, expectedGrid, response["gameGrid"])
}

func addPlayerTwoToGameSessionIntegration(t *testing.T) {
	jsonBody := `{"player2": "Alice"}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId+"/players", "application/json", bodyReader)
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

	assert.Equal(t, "John", response["player1"])
	assert.Equal(t, "Alice", response["player2"])
	assert.Equal(t, "John", response["nextPlayerMove"])

	expectedGrid := []interface{}{
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
	}
	assert.Equal(t, expectedGrid, response["gameGrid"])
}
