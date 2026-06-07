package detector

import (
	"fmt"

	"github.com/dmcelhill/stream-intel/pkg/aggregator"
)

type FareSpikeDetector struct {
	Threshold float64
}

func (d *FareSpikeDetector) Detect(snapshots map[int]aggregator.ZoneSnapshot) []Alert {
	alerts := []Alert{}
	for zone, snapshot := range snapshots {
		if snapshot.AvgFare > d.Threshold {
			alerts = append(alerts, Alert{
				Zone:    zone,
				Type:    "FareSpike",
				Message: fmt.Sprintf("Average fare in zone %d has exceeded threshold $%.2f", zone, d.Threshold),
			})
		}
	}
	return alerts
}
