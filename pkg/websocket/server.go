package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ClientMessage struct {
	Action string `json:"action"`
}

type ReplayFunc func() error

type Server struct {
	hub      *Hub
	logger   *slog.Logger
	onReplay ReplayFunc
}

func NewServer(hub *Hub, logger *slog.Logger) *Server {
	return &Server{hub: hub, logger: logger}
}

func (s *Server) SetReplayFunc(fn ReplayFunc) {
	s.onReplay = fn
}

func (s *Server) HandleConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("WebSocket upgrade failed", "error", err)
		return
	}

	client := &Client{send: make(chan []byte, 64)}
	s.hub.Register(client)

	// Write goroutine: reads from client.send, writes to WebSocket
	go func() {
		defer conn.Close()
		defer s.hub.Unregister(client)
		for msg := range client.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()

	// Read goroutine: handles client commands and detects disconnect
	go func() {
		defer s.hub.Unregister(client)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var cmd ClientMessage
			if err := json.Unmarshal(msg, &cmd); err != nil {
				continue
			}
			if cmd.Action == "replay" && s.onReplay != nil {
				if err := s.onReplay(); err != nil {
					s.logger.Error("Replay failed", "error", err)
				}
			}
		}
	}()
}

func (s *Server) Run(ctx context.Context, port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.HandleConnect)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}
