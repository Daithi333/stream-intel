package websocket

import (
	"sync"
	"testing"
	"time"
)

func TestRegisterAndBroadcast(t *testing.T) {
	hub := NewHub()
	client := &Client{send: make(chan []byte, 10)}
	hub.Register(client)

	hub.Broadcast([]byte("hello"))

	select {
	case msg := <-client.send:
		if string(msg) != "hello" {
			t.Errorf("expected 'hello', got '%s'", string(msg))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for broadcast")
	}
}

func TestUnregisterClosesChannel(t *testing.T) {
	hub := NewHub()
	client := &Client{send: make(chan []byte, 10)}
	hub.Register(client)
	hub.Unregister(client)

	_, open := <-client.send
	if open {
		t.Error("expected client channel to be closed after unregister")
	}
}

func TestUnregisterIdempotent(t *testing.T) {
	hub := NewHub()
	client := &Client{send: make(chan []byte, 10)}
	hub.Register(client)
	hub.Unregister(client)
	hub.Unregister(client) // should not panic
}

func TestBroadcastToMultipleClients(t *testing.T) {
	hub := NewHub()
	c1 := &Client{send: make(chan []byte, 10)}
	c2 := &Client{send: make(chan []byte, 10)}
	hub.Register(c1)
	hub.Register(c2)

	hub.Broadcast([]byte("data"))

	for _, c := range []*Client{c1, c2} {
		select {
		case msg := <-c.send:
			if string(msg) != "data" {
				t.Errorf("expected 'data', got '%s'", string(msg))
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for broadcast to client")
		}
	}
}

func TestBroadcastDropsWhenBufferFull(t *testing.T) {
	dropped := 0
	hub := NewHub()
	hub.OnDrop = func() { dropped++ }

	client := &Client{send: make(chan []byte, 1)} // buffer of 1
	hub.Register(client)

	hub.Broadcast([]byte("first"))  // fills buffer
	hub.Broadcast([]byte("second")) // dropped

	if dropped != 1 {
		t.Errorf("expected 1 drop, got %d", dropped)
	}

	msg := <-client.send
	if string(msg) != "first" {
		t.Errorf("expected 'first', got '%s'", string(msg))
	}
}

func TestBroadcastConcurrentSafety(t *testing.T) {
	hub := NewHub()
	client := &Client{send: make(chan []byte, 100)}
	hub.Register(client)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hub.Broadcast([]byte("msg"))
		}()
	}
	wg.Wait()

	count := len(client.send)
	if count != 50 {
		t.Errorf("expected 50 messages, got %d", count)
	}
}
