package websocket

import (
	"sync"
)

type Client struct {
	send chan []byte
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
	OnDrop  func() // callback for when a message is dropped
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = true
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		close(c.send)
		delete(h.clients, c)
	}
}

func (h *Hub) Broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			// Client's buffer is full so drop and record metric
			if h.OnDrop != nil {
				h.OnDrop()
			}
		}
	}
}
