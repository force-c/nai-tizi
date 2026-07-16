package websocket

import (
	"net/http/httptest"
	"strings"
	"testing"

	logging "github.com/gcc798/quick.admin/internal/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

func TestHandlerUpgradesEchoResponse(t *testing.T) {
	log, err := logging.NewLoggerWithConfig(&logging.Config{Level: "error", Output: "console", Encoding: "json"})
	require.NoError(t, err)
	hub := NewHub(log)
	require.NoError(t, hub.Start())
	t.Cleanup(func() { _ = hub.Stop() })

	e := echo.New()
	handler := NewHandler(hub, log)
	e.GET("/resource/websocket", func(c *echo.Context) error {
		handler.ServeWs(c)
		return nil
	})
	server := httptest.NewServer(e)
	t.Cleanup(server.Close)

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/resource/websocket?userId=42"
	conn, response, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	require.Equal(t, 101, response.StatusCode)
	require.NoError(t, conn.Close())
}
