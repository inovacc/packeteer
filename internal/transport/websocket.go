package transport

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
)

// WebSocketProxy bridges WebSocket connections to an HTTP MCP endpoint.
// WebSocket clients send JSON-RPC requests as text frames; the proxy
// forwards them to the Streamable HTTP handler and relays responses back.
type WebSocketProxy struct {
	upstream http.Handler
	logger   *slog.Logger
	mu       sync.Mutex
	active   int
}

// NewWebSocketProxy creates a proxy that accepts WebSocket connections
// and forwards JSON-RPC messages to the given HTTP handler.
func NewWebSocketProxy(upstream http.Handler, logger *slog.Logger) *WebSocketProxy {
	return &WebSocketProxy{
		upstream: upstream,
		logger:   logger,
	}
}

// ServeHTTP upgrades the connection to WebSocket and proxies messages.
func (p *WebSocketProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		p.logger.Error("websocket accept failed", "error", err)
		return
	}
	defer conn.CloseNow()

	p.mu.Lock()
	p.active++
	sessionID := p.active
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		p.active--
		p.mu.Unlock()
	}()

	p.logger.Info("websocket client connected", "session", sessionID, "remote", r.RemoteAddr)
	defer p.logger.Info("websocket client disconnected", "session", sessionID)

	ctx := r.Context()

	for {
		_, msg, err := conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				return
			}
			if ctx.Err() != nil {
				return
			}
			p.logger.Debug("websocket read error", "error", err, "session", sessionID)
			return
		}

		// Forward the JSON-RPC request to the upstream HTTP handler.
		resp, err := forwardToHTTP(ctx, p.upstream, msg, r)
		if err != nil {
			p.logger.Error("upstream forward failed", "error", err, "session", sessionID)
			continue
		}

		if err := conn.Write(ctx, websocket.MessageText, resp); err != nil {
			p.logger.Debug("websocket write error", "error", err, "session", sessionID)
			return
		}
	}
}

// forwardToHTTP sends a JSON-RPC message to the upstream handler via an in-process HTTP roundtrip.
func forwardToHTTP(ctx context.Context, handler http.Handler, body []byte, originalReq *http.Request) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "/", io.NopCloser(
		&bytesReader{data: body, pos: 0},
	))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Copy session header if present.
	if sid := originalReq.Header.Get("Mcp-Session-Id"); sid != "" {
		req.Header.Set("Mcp-Session-Id", sid)
	}

	rec := &responseRecorder{headers: make(http.Header)}
	handler.ServeHTTP(rec, req)

	return rec.body, nil
}

type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

type responseRecorder struct {
	statusCode int
	headers    http.Header
	body       []byte
}

func (r *responseRecorder) Header() http.Header         { return r.headers }
func (r *responseRecorder) WriteHeader(statusCode int)   { r.statusCode = statusCode }
func (r *responseRecorder) Write(b []byte) (int, error)  { r.body = append(r.body, b...); return len(b), nil }

// ActiveConnections returns the number of active WebSocket sessions.
func (p *WebSocketProxy) ActiveConnections() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.active
}
