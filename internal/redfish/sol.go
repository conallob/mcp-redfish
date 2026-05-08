package redfish

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	solLineCapacity  = 1000
	solReconnectWait = 5 * time.Second
)

// SOLState represents the connection state of the SOL client.
type SOLState string

const (
	SOLStateDisconnected SOLState = "disconnected"
	SOLStateConnecting   SOLState = "connecting"
	SOLStateConnected    SOLState = "connected"
	SOLStateError        SOLState = "error"
)

// SOLClient manages a WebSocket Serial-over-LAN connection and buffers console output.
type SOLClient struct {
	wsURL    string
	username string
	password string
	insecure bool

	buf lineBuffer

	mu      sync.RWMutex
	state   SOLState
	lastErr string

	cancel context.CancelFunc
}

// NewPreloadedSOLClient creates a connected SOLClient with the given lines already buffered.
// Intended for testing and dependency injection scenarios.
func NewPreloadedSOLClient(lines []string) *SOLClient {
	s := &SOLClient{
		buf:   newLineBuffer(solLineCapacity),
		state: SOLStateConnected,
	}
	for _, l := range lines {
		s.buf.append(l)
	}
	return s
}

// NewSOLClient creates an SOLClient that connects to wsURL using credentials from cfg.
func NewSOLClient(cfg Config, wsURL string) *SOLClient {
	return &SOLClient{
		wsURL:    wsURL,
		username: cfg.Username,
		password: cfg.Password,
		insecure: cfg.Insecure,
		buf:      newLineBuffer(solLineCapacity),
		state:    SOLStateDisconnected,
	}
}

// Start begins the background WebSocket connection and reconnection loop.
func (s *SOLClient) Start(ctx context.Context) {
	ctx, s.cancel = context.WithCancel(ctx)
	go s.loop(ctx)
}

// Stop closes the SOL connection and terminates the background loop.
func (s *SOLClient) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// State returns the current connection state and last error message (if any).
func (s *SOLClient) State() (SOLState, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state, s.lastErr
}

// RecentLines returns up to n of the most recently buffered console lines.
// If n <= 0, all buffered lines are returned.
func (s *SOLClient) RecentLines(n int) []string {
	count := s.buf.count()
	if n <= 0 || n > count {
		n = count
	}
	return s.buf.recent(n)
}

func (s *SOLClient) setState(state SOLState, errMsg string) {
	s.mu.Lock()
	s.state = state
	s.lastErr = errMsg
	s.mu.Unlock()
}

func (s *SOLClient) loop(ctx context.Context) {
	for {
		s.setState(SOLStateConnecting, "")
		err := s.connect(ctx)
		if ctx.Err() != nil {
			s.setState(SOLStateDisconnected, "")
			return
		}
		if err != nil {
			s.setState(SOLStateError, err.Error())
		}
		select {
		case <-ctx.Done():
			s.setState(SOLStateDisconnected, "")
			return
		case <-time.After(solReconnectWait):
		}
	}
}

func (s *SOLClient) connect(ctx context.Context) error {
	dialer := websocket.Dialer{
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: s.insecure}, //nolint:gosec
		HandshakeTimeout: 15 * time.Second,
	}

	headers := make(http.Header)
	if s.username != "" {
		creds := base64.StdEncoding.EncodeToString([]byte(s.username + ":" + s.password))
		headers.Set("Authorization", "Basic "+creds)
	}

	conn, resp, err := dialer.DialContext(ctx, s.wsURL, headers)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil {
		return fmt.Errorf("dial %s: %w", s.wsURL, err)
	}
	defer func() { _ = conn.Close() }()

	s.setState(SOLStateConnected, "")

	// Close the connection when the context is cancelled.
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-done:
		}
	}()

	partial := ""
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		// Append to any partial line carried over from the previous frame, then
		// split on newlines. The last segment may be an incomplete line.
		text := partial + string(data)
		parts := strings.Split(text, "\n")
		for i, part := range parts {
			if i < len(parts)-1 {
				s.buf.append(strings.TrimRight(part, "\r"))
			}
		}
		partial = parts[len(parts)-1]
	}
}

// lineBuffer is a fixed-capacity circular buffer of strings.
type lineBuffer struct {
	mu    sync.RWMutex
	lines []string
	head  int
	total int
}

func newLineBuffer(capacity int) lineBuffer {
	return lineBuffer{lines: make([]string, capacity)}
}

func (b *lineBuffer) append(line string) {
	b.mu.Lock()
	b.lines[b.head] = line
	b.head = (b.head + 1) % len(b.lines)
	b.total++
	b.mu.Unlock()
}

func (b *lineBuffer) count() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.total < len(b.lines) {
		return b.total
	}
	return len(b.lines)
}

func (b *lineBuffer) recent(n int) []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	stored := b.total
	if stored > len(b.lines) {
		stored = len(b.lines)
	}
	if n > stored {
		n = stored
	}
	if n == 0 {
		return nil
	}
	result := make([]string, n)
	start := (b.head - n + len(b.lines)) % len(b.lines)
	for i := 0; i < n; i++ {
		result[i] = b.lines[(start+i)%len(b.lines)]
	}
	return result
}
