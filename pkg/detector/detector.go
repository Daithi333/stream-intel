package detector

import "github.com/dmcelhill/stream-intel/pkg/aggregator"

type Alert struct {
	Zone    int
	Type    string
	Message string
}

type Detector interface {
	Detect(snapshots map[int]aggregator.ZoneSnapshot) []Alert
}
