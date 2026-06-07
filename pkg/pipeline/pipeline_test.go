package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/dmcelhill/stream-intel/pkg/model"
)

func TestSendAndReceive(t *testing.T) {
	p := New(10)
	defer p.Close()

	ctx := context.Background()
	event := model.TaxiTrip{
		EventId:      "test-1",
		FareAmount:   25.50,
		PuLocationId: 42,
	}

	if err := p.Send(ctx, event); err != nil {
		t.Fatalf("unexpected error on Send: %v", err)
	}

	select {
	case received := <-p.Receive():
		if received.EventId != "test-1" {
			t.Errorf("expected EventId test-1, got %s", received.EventId)
		}
		if received.FareAmount != 25.50 {
			t.Errorf("expected FareAmount 25.50, got %f", received.FareAmount)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestSendReturnsErrorOnCancelledContext(t *testing.T) {
	p := New(0) // unbuffered — Send will block unless someone is receiving
	defer p.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	event := model.TaxiTrip{EventId: "test-2"}
	err := p.Send(ctx, event)

	if err == nil {
		t.Fatal("expected error on cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestReceiveExitsOnClose(t *testing.T) {
	p := New(5)

	ctx := context.Background()
	_ = p.Send(ctx, model.TaxiTrip{EventId: "a"})
	_ = p.Send(ctx, model.TaxiTrip{EventId: "b"})
	p.Close()

	var received []string
	for trip := range p.Receive() {
		received = append(received, trip.EventId)
	}

	if len(received) != 2 {
		t.Errorf("expected 2 events, got %d", len(received))
	}
}

func TestBackpressure(t *testing.T) {
	p := New(1) // buffer of 1
	defer p.Close()

	ctx := context.Background()
	_ = p.Send(ctx, model.TaxiTrip{EventId: "fills-buffer"})

	// Second send should block — use a timeout context to prove it
	timeoutCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	err := p.Send(timeoutCtx, model.TaxiTrip{EventId: "blocked"})
	if err == nil {
		t.Fatal("expected timeout error when buffer is full, got nil")
	}
}
