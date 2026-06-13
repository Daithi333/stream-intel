package sink

import "github.com/dmcelhill/stream-intel/pkg/detector"

type Sink interface {
	Send(alert detector.Alert) error
}
