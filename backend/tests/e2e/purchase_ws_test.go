package tests_e2e

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestGetPurchaseWebsocketWithInvalidToken(t *testing.T) {
	ts, cleanup := setupTestEnvironment(t)
	defer cleanup()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/api/v2/purchases/123/ws"

	// no token provided
	_, resp, err := connectWS(t, wsURL, "", "http://localhost:3000")
	require.Error(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// with invalid token
	_, resp, err = connectWS(t, wsURL, "invalid.token.here", "http://localhost:3000")
	require.Error(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetPurchaseWebsocketWithInvalidOrigin(t *testing.T) {
	ts, cleanup := setupTestEnvironment(t)
	defer cleanup()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/api/v2/purchases/01982971-a954-74ed-9735-a75e08efa8f6/ws"
	token := getJwtForDemoUser()

	conn, resp, err := connectWS(t, wsURL, token, "http://example.com:3000")
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusForbidden, resp.StatusCode) // oder 400, je nach Upgrader-Fehler

	if conn != nil {
		conn.Close()
	}
}

func TestGetPurchaseWebsocketWithValidToken(t *testing.T) {
	ts, cleanup := setupTestEnvironment(t)
	defer cleanup()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/api/v2/purchases/01982971-a954-74ed-9735-a75e08efa8f6/ws"
	token := getJwtForDemoUser()

	conn, resp, err := connectWS(t, wsURL, token, "http://localhost:3000")
	require.NoError(t, err)
	require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

	defer conn.Close()
}

func connectWS(t *testing.T, url string, token string, origin string) (*websocket.Conn, *http.Response, error) {
	t.Helper()

	dialer := websocket.Dialer{
		Subprotocols: []string{token},
	}

	reqHeader := http.Header{}
	reqHeader.Set("Origin", origin)

	conn, resp, err := dialer.Dial(url, reqHeader)

	return conn, resp, err
}
