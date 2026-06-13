package sink

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/dmcelhill/stream-intel/pkg/detector"
)

func TestLogSinkSendsAlert(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	s := &LogSink{Logger: logger}

	alert := detector.Alert{
		Zone:    42,
		Type:    "FareSpike",
		Message: "test alert",
	}

	err := s.Send(alert)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Fatal("expected log output, got empty string")
	}
	if !bytes.Contains([]byte(output), []byte("test alert")) {
		t.Errorf("expected output to contain 'test alert', got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("42")) {
		t.Errorf("expected output to contain zone '42', got: %s", output)
	}
}

func TestLogSinkSatisfiesSinkInterface(t *testing.T) {
	var s Sink = &LogSink{Logger: slog.Default()}
	_ = s // compiles means it satisfies the interface
}
