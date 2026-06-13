package sink

import (
	"log/slog"

	"github.com/dmcelhill/stream-intel/pkg/detector"
)

type LogSink struct {
	Logger *slog.Logger
}

func (s *LogSink) Send(alert detector.Alert) error {
	s.Logger.Info("Alert",
		"zone", alert.Zone,
		"type", alert.Type,
		"message", alert.Message,
	)
	return nil
}
