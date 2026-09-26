package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type pingMock struct {
	sqlite.RepositoryInterface

	err error
}

func (m pingMock) Ping(context.Context) error {
	return m.err
}

func TestGetHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", http.NoBody)

	handler.GetHealth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ok", body["status"])
}

func TestGetReady_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{repo: pingMock{err: nil}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ready", http.NoBody)

	handler.GetReady(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ready", body["status"])
}

func TestGetReady_DatabaseDown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{repo: pingMock{err: errors.New("db down")}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ready", http.NoBody)

	handler.GetReady(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "unavailable", body["status"])
	assert.Equal(t, "database down", body["error"])
}

// The route is unauthenticated, so the response must stay generic no matter how
// specific the underlying failure is, and must not be written twice.
func TestGetReady_DoesNotLeakInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &Handler{repo: pingMock{err: errors.New(
		"failed to ping database: failed to open data/kasseapparat.db: disk I/O error",
	)}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ready", http.NoBody)

	handler.GetReady(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, `{"error":"database down","status":"unavailable"}`, w.Body.String())
}
