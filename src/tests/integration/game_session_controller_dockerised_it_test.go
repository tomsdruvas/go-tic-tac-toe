package integration_test

import (
	"encoding/json"
	"fmt"
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
	playerOneMakesFirstMove(t)
	playerTwoMakesSecondMove(t)
	playerOneMakesThirdMove(t)
	playerTwoMakesFourthMove(t)
	playerOneMakesFifthMove(t)
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
		[]interface{}{"Empty", "Empty", "Empty"},
		[]interface{}{"Empty", "Empty", "Empty"},
		[]interface{}{"Empty", "Empty", "Empty"},
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{
		"sessionId": "%s",
		"player1": "John",
		"nextPlayerMove": "John",
		"gameSessionStatus": "Active",
		"gameGrid": [
			["Empty", "Empty", "Empty"],
			["Empty", "Empty", "Empty"],
			["Empty", "Empty", "Empty"]
		]
	}`, currentSessionId)

	assert.JSONEq(t, expectedJSON, string(body))
}

func assertBodyWithTwoPlayer(t *testing.T, resp *http.Response, err error) {
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{
		"sessionId": "%s",
		"player1": "John",
		"player2": "Alice",
		"nextPlayerMove": "John",
		"gameSessionStatus": "Active",
		"gameGrid": [
			["Empty", "Empty", "Empty"],
			["Empty", "Empty", "Empty"],
			["Empty", "Empty", "Empty"]
		]
	}`, currentSessionId)

	assert.JSONEq(t, expectedJSON, string(body))
}

func addPlayerTwoToGameSessionIntegration(t *testing.T) {
	jsonBody := `{"player2": "Alice"}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId+"/players", "application/json", bodyReader)
	assertBodyWithTwoPlayer(t, resp, err)
}

func playerOneMakesFirstMove(t *testing.T) {
	jsonBody := `{"playerName": "John",
					"xAxis": 0, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId, "application/json", bodyReader)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{
		"sessionId": "%s",
		"player1": "John",
		"player2": "Alice",
		"nextPlayerMove": "Alice",
		"gameSessionStatus": "Active",
		"gameGrid": [
			["Cross", "Empty", "Empty"],
			["Empty", "Empty", "Empty"],
			["Empty", "Empty", "Empty"]
		]
	}`, currentSessionId)

	assert.JSONEq(t, expectedJSON, string(body))
}

func playerTwoMakesSecondMove(t *testing.T) {
	jsonBody := `{"playerName": "Alice",
					"xAxis": 1, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId, "application/json", bodyReader)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{
		"sessionId": "%s",
		"player1": "John",
		"player2": "Alice",
		"nextPlayerMove": "John",
		"gameSessionStatus": "Active",
		"gameGrid": [
			["Cross", "Empty", "Empty"],
			["Circle", "Empty", "Empty"],
			["Empty", "Empty", "Empty"]
		]
	}`, currentSessionId)

	assert.JSONEq(t, expectedJSON, string(body))
}

func playerOneMakesThirdMove(t *testing.T) {
	jsonBody := `{"playerName": "John",
					"xAxis": 0, "yAxis": 1}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId, "application/json", bodyReader)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{
		"sessionId": "%s",
		"player1": "John",
		"player2": "Alice",
		"nextPlayerMove": "Alice",
		"gameSessionStatus": "Active",
		"gameGrid": [
			["Cross", "Cross", "Empty"],
			["Circle", "Empty", "Empty"],
			["Empty", "Empty", "Empty"]
		]
	}`, currentSessionId)

	assert.JSONEq(t, expectedJSON, string(body))
}

func playerTwoMakesFourthMove(t *testing.T) {
	jsonBody := `{"playerName": "Alice",
					"xAxis": 2, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId, "application/json", bodyReader)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{
		"sessionId": "%s",
		"player1": "John",
		"player2": "Alice",
		"nextPlayerMove": "John",
		"gameSessionStatus": "Active",
		"gameGrid": [
			["Cross", "Cross", "Empty"],
			["Circle", "Empty", "Empty"],
			["Circle", "Empty", "Empty"]
		]
	}`, currentSessionId)

	assert.JSONEq(t, expectedJSON, string(body))
}

func playerOneMakesFifthMove(t *testing.T) {
	jsonBody := `{"playerName": "John",
					"xAxis": 0, "yAxis": 2}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId, "application/json", bodyReader)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{
		"sessionId": "%s",
		"player1": "John",
		"player2": "Alice",
		"gameSessionStatus": "Finished",
		"winner": "John",
		"gameGrid": [
			["Cross", "Cross", "Cross"],
			["Circle", "Empty", "Empty"],
			["Circle", "Empty", "Empty"]
		]
	}`, currentSessionId)

	assert.JSONEq(t, expectedJSON, string(body))
}
