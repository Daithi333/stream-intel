package websocket

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func setupTestServer(t *testing.T) (*Server, *httptest.Server) {
	t.Helper()
	hub := NewHub()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := NewServer(hub, logger)

	ts := httptest.NewServer(http.HandlerFunc(srv.HandleConnect))
	return srv, ts
}

func dialWS(t *testing.T, ts *httptest.Server) *websocket.Conn {
	t.Helper()
	url := "ws" + strings.TrimPrefix(ts.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("failed to dial websocket: %v", err)
	}
	return conn
}

func TestClientReceivesBroadcast(t *testing.T) {
	srv, ts := setupTestServer(t)
	defer ts.Close()

	conn := dialWS(t, ts)
	defer conn.Close()

	time.Sleep(50 * time.Millisecond) // allow registration

	srv.hub.Broadcast([]byte(`{"zone":1}`))

	_ = conn.SetReadDeadline(time.Now().Add(time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}
	if string(msg) != `{"zone":1}` {
		t.Errorf("expected '{\"zone\":1}', got '%s'", string(msg))
	}
}

func TestReplayCommand(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	replayCalled := false
	srv, _ := setupTestServer(t)
	srv.SetReplayFunc(func() error {
		replayCalled = true
		return nil
	})

	ts2 := httptest.NewServer(http.HandlerFunc(srv.HandleConnect))
	defer ts2.Close()

	conn := dialWS(t, ts2)
	defer conn.Close()

	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"action":"replay"}`))
	if err != nil {
		t.Fatalf("failed to send replay command: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if !replayCalled {
		t.Error("expected replay function to be called")
	}
}

func TestMultipleClientsReceiveBroadcast(t *testing.T) {
	srv, ts := setupTestServer(t)
	defer ts.Close()

	conn1 := dialWS(t, ts)
	defer conn1.Close()
	conn2 := dialWS(t, ts)
	defer conn2.Close()

	time.Sleep(50 * time.Millisecond)

	srv.hub.Broadcast([]byte("update"))

	for _, conn := range []*websocket.Conn{conn1, conn2} {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}
		if string(msg) != "update" {
			t.Errorf("expected 'update', got '%s'", string(msg))
		}
	}
}

func TestDisconnectUnregistersClient(t *testing.T) {
	srv, ts := setupTestServer(t)
	defer ts.Close()

	conn := dialWS(t, ts)
	time.Sleep(50 * time.Millisecond)

	srv.hub.mu.RLock()
	countBefore := len(srv.hub.clients)
	srv.hub.mu.RUnlock()

	conn.Close()
	time.Sleep(100 * time.Millisecond)

	srv.hub.mu.RLock()
	countAfter := len(srv.hub.clients)
	srv.hub.mu.RUnlock()

	if countBefore != 1 {
		t.Errorf("expected 1 client before disconnect, got %d", countBefore)
	}
	if countAfter != 0 {
		t.Errorf("expected 0 clients after disconnect, got %d", countAfter)
	}
}
