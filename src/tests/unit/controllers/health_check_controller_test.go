package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"tic-tac-toe-game/src/controllers"
)

func setupHealthCheckRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	hc := controllers.NewHealthCheckController()
	r.GET("/health", hc.HealthCheckHandler)

	return r
}

func TestHealthCheckHandler(t *testing.T) {
	r := setupHealthCheckRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "healthy", "message": "Service is up and running"}`, w.Body.String())
}
