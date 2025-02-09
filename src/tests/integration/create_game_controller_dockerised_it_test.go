package controllers_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCreateGameIntegration(t *testing.T) {
	jsonBody := `{"player1": "John"}`
	bodyReader := strings.NewReader(jsonBody)
	resp, err := http.Post(TestServerURL+"/create-game", "application/json", bodyReader)
	assert.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	sessionId, ok := response["sessionId"].(string)
	assert.True(t, ok, "sessionId should be a string")
	assert.NotEmpty(t, sessionId, "sessionId should not be empty")

	assert.Equal(t, "John", response["player1"])
	assert.Equal(t, "John", response["nextPlayerMove"])

	expectedGrid := []interface{}{
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
		[]interface{}{float64(0), float64(0), float64(0)},
	}
	assert.Equal(t, expectedGrid, response["gameGrid"])
}
