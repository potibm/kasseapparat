package tests_e2e

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

const originURL = "http://localhost:3000"

func TestGetPurchaseWebsocketWithInvalidOrigin(t *testing.T) {
	ts, cleanup := setupTestEnvironment(t)
	defer cleanup()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/api/v3/purchases/01982971-a954-74ed-9735-a75e08efa8f6/ws"

	conn, resp, err := connectWS(t, wsURL, "http://example.com:3000")
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	if conn != nil {
		conn.Close()
	}
}

func TestGetPurchaseWebsocketWithValidOrigin(t *testing.T) {
	ts, cleanup := setupTestEnvironment(t)
	defer cleanup()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/api/v3/purchases/01982971-a954-74ed-9735-a75e08efa8f6/ws"

	conn, resp, err := connectWS(t, wsURL, originURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

	defer conn.Close()
}

func connectWS(t *testing.T, url, origin string) (*websocket.Conn, *http.Response, error) {
	t.Helper()

	dialer := websocket.Dialer{}

	reqHeader := http.Header{}
	reqHeader.Set("Origin", origin)
	reqHeader.Set("X-Remote-User", "testuser")

	return dialer.Dial(url, reqHeader)
}
