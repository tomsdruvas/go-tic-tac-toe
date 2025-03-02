package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	_ "tic-tac-toe-game/src/tests/integration/setup"
)

const (
	websocketGameSessionPath = "ws://localhost:8091/game-session?gameSessionId="
)

var (
	currentSessionId string
	conn             *websocket.Conn
	httpResponse     *http.Response
	err              error
	done             chan struct{}
	messages         []string
	mu               sync.Mutex
)

func TestTicTacToeFullFlow(t *testing.T) {
	createGameSessionIntegration(t)

	connectToWebsocketsForGameSession(t)

	assertBodyWithOnePlayer(retrieveGameSessionWithResponseBodyIntegration(t))
	addPlayerTwoToGameSessionIntegration(t)
	assertBodyWithTwoPlayer(retrieveGameSessionWithResponseBodyIntegration(t))
	playerOneMakesFirstMove(t)
	playerTwoMakesSecondMove(t)
	playerOneMakesThirdMove(t)
	playerTwoMakesFourthMove(t)
	playerOneMakesFifthMove(t)

	gracefullyCloseWebSocketConnection(t)
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

func connectToWebsocketsForGameSession(t *testing.T) {
	dialer := websocket.Dialer{}
	conn, httpResponse, err = dialer.Dial(websocketGameSessionPath+currentSessionId, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	if conn == nil {
		t.Fatalf("Expected a valid WebSocket connection, got nil")
	}

	if httpResponse.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("Expected HTTP status 101 (Switching Protocols), got %v", httpResponse.StatusCode)
	}

	done = make(chan struct{})

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				select {
				case <-done:
					t.Logf("Gracefully stopping message read loop.")
					return
				default:
					t.Logf("Error reading message: %v", err)
					return
				}
			}

			t.Logf("Received message: %s", msg)
			mu.Lock()
			messages = append(messages, string(msg))
			mu.Unlock()
		}
	}()
}

func gracefullyCloseWebSocketConnection(t *testing.T) {
	close(done)
	<-done
	err = conn.Close()
	if err != nil {
		t.Fatalf("Error closing WebSocket connection: %v", err)
	}
	t.Log("WebSocket connection closed gracefully")
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
	assertAndPopWebsocketMessage(t, expectedJSON)
}

func assertAndPopWebsocketMessage(t *testing.T, expectedMessage string) {
	mu.Lock()
	defer mu.Unlock()

	assert.Equal(t, 1, len(messages))
	assert.JSONEq(t, messages[0], expectedMessage)

	messages = messages[1:]
}

func playerOneMakesFirstMove(t *testing.T) {
	jsonBody := `{"playerName": "John",
					"xAxis": 0, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId+"/move", "application/json", bodyReader)
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
	assertAndPopWebsocketMessage(t, expectedJSON)
}

func playerTwoMakesSecondMove(t *testing.T) {
	jsonBody := `{"playerName": "Alice",
					"xAxis": 1, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId+"/move", "application/json", bodyReader)
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
	assertAndPopWebsocketMessage(t, expectedJSON)
}

func playerOneMakesThirdMove(t *testing.T) {
	jsonBody := `{"playerName": "John",
					"xAxis": 0, "yAxis": 1}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId+"/move", "application/json", bodyReader)
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
	assertAndPopWebsocketMessage(t, expectedJSON)
}

func playerTwoMakesFourthMove(t *testing.T) {
	jsonBody := `{"playerName": "Alice",
					"xAxis": 2, "yAxis": 0}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId+"/move", "application/json", bodyReader)
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
	assertAndPopWebsocketMessage(t, expectedJSON)
}

func playerOneMakesFifthMove(t *testing.T) {
	jsonBody := `{"playerName": "John",
					"xAxis": 0, "yAxis": 2}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/game-session/"+currentSessionId+"/move", "application/json", bodyReader)
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
	assertAndPopWebsocketMessage(t, expectedJSON)
}
