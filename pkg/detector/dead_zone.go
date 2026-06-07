package detector

import (
	"fmt"
	"time"

	"github.com/dmcelhill/stream-intel/pkg/aggregator"
)

type DeadZoneDetector struct {
	InactiveThreshold time.Duration
}

func (d *DeadZoneDetector) Detect(snapshots map[int]aggregator.ZoneSnapshot) []Alert {
	alerts := []Alert{}
	for zone, snapshot := range snapshots {
		if time.Since(snapshot.LastSeen) > d.InactiveThreshold {
			alerts = append(alerts, Alert{
				Zone:    zone,
				Type:    "DeadZone",
				Message: fmt.Sprintf("Zone %d has not had activity since %s", zone, d.InactiveThreshold),
			})
		}
	}
	return alerts
}
