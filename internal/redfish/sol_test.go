package redfish_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/conallob/mcp-redfish/internal/redfish"
)

// --- lineBuffer tests (via SOLClient) ---

func TestSOL_RecentLines_Empty(t *testing.T) {
	sol := redfish.NewSOLClient(redfish.Config{}, "ws://unused")
	if got := sol.RecentLines(10); len(got) != 0 {
		t.Errorf("expected 0 lines, got %d", len(got))
	}
}

func TestSOL_RecentLines_BelowCapacity(t *testing.T) {
	srv := newSOLTestServer(t, []string{"line1", "line2", "line3"})
	sol := redfish.NewSOLClient(redfish.Config{}, srv.wsURL)
	sol.Start(context.Background())
	defer sol.Stop()

	waitForLines(t, sol, 3, 2*time.Second)

	lines := sol.RecentLines(10)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line1" || lines[1] != "line2" || lines[2] != "line3" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestSOL_RecentLines_LimitRespected(t *testing.T) {
	srv := newSOLTestServer(t, []string{"a", "b", "c", "d", "e"})
	sol := redfish.NewSOLClient(redfish.Config{}, srv.wsURL)
	sol.Start(context.Background())
	defer sol.Stop()

	waitForLines(t, sol, 5, 2*time.Second)

	lines := sol.RecentLines(3)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// Should be the last 3: c, d, e
	if lines[0] != "c" || lines[1] != "d" || lines[2] != "e" {
		t.Errorf("expected last 3 lines [c d e], got %v", lines)
	}
}

func TestSOL_StateTransitions(t *testing.T) {
	srv := newSOLTestServer(t, nil)
	sol := redfish.NewSOLClient(redfish.Config{}, srv.wsURL)

	state, _ := sol.State()
	if state != redfish.SOLStateDisconnected {
		t.Errorf("expected disconnected before start, got %s", state)
	}

	sol.Start(context.Background())
	defer sol.Stop()

	// Wait for connected state.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if s, _ := sol.State(); s == redfish.SOLStateConnected {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Errorf("SOL client did not reach connected state within timeout")
}

func TestSOL_StopSetsDisconnected(t *testing.T) {
	srv := newSOLTestServer(t, nil)
	sol := redfish.NewSOLClient(redfish.Config{}, srv.wsURL)
	sol.Start(context.Background())

	waitForState(t, sol, redfish.SOLStateConnected, 2*time.Second)

	sol.Stop()
	waitForState(t, sol, redfish.SOLStateDisconnected, 2*time.Second)
}

func TestSOL_ErrorOnBadURL(t *testing.T) {
	sol := redfish.NewSOLClient(redfish.Config{}, "ws://127.0.0.1:1") // port 1 should refuse
	sol.Start(context.Background())
	defer sol.Stop()

	waitForState(t, sol, redfish.SOLStateError, 3*time.Second)
	_, errMsg := sol.State()
	if errMsg == "" {
		t.Error("expected a non-empty error message")
	}
}

// --- helpers ---

var upgrader = websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}

type solTestServer struct {
	wsURL string
}

// newSOLTestServer starts a test WebSocket server that sends lines to each connecting client.
func newSOLTestServer(t *testing.T, lines []string) *solTestServer {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/sol", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		if len(lines) > 0 {
			payload := strings.Join(lines, "\n") + "\n"
			_ = conn.WriteMessage(websocket.TextMessage, []byte(payload))
		}
		// Keep connection open until client disconnects.
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/sol"
	return &solTestServer{wsURL: wsURL}
}

func waitForLines(t *testing.T, sol *redfish.SOLClient, want int, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if len(sol.RecentLines(0)) >= want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %d buffered lines (got %d)", want, len(sol.RecentLines(0)))
}

func waitForState(t *testing.T, sol *redfish.SOLClient, want redfish.SOLState, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if s, _ := sol.State(); s == want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	s, errMsg := sol.State()
	t.Fatalf("timed out waiting for state %s (current: %s, err: %s)", want, s, errMsg)
}
