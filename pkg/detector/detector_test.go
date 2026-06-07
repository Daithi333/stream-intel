package detector

import (
	"testing"
	"time"

	"github.com/dmcelhill/stream-intel/pkg/aggregator"
)

func TestFareSpikeDetectsHighFare(t *testing.T) {
	d := &FareSpikeDetector{Threshold: 50.0}
	snapshots := map[int]aggregator.ZoneSnapshot{
		1: {TripCount: 10, AvgFare: 60.0, LastSeen: time.Now()},
		2: {TripCount: 5, AvgFare: 30.0, LastSeen: time.Now()},
	}

	alerts := d.Detect(snapshots)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Zone != 1 {
		t.Errorf("expected alert for zone 1, got zone %d", alerts[0].Zone)
	}
	if alerts[0].Type != "FareSpike" {
		t.Errorf("expected type FareSpike, got %s", alerts[0].Type)
	}
}

func TestFareSpikeNoAlertBelowThreshold(t *testing.T) {
	d := &FareSpikeDetector{Threshold: 50.0}
	snapshots := map[int]aggregator.ZoneSnapshot{
		1: {TripCount: 10, AvgFare: 25.0, LastSeen: time.Now()},
		2: {TripCount: 5, AvgFare: 49.99, LastSeen: time.Now()},
	}

	alerts := d.Detect(snapshots)

	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestDeadZoneDetectsInactiveZone(t *testing.T) {
	d := &DeadZoneDetector{InactiveThreshold: 5 * time.Minute}
	snapshots := map[int]aggregator.ZoneSnapshot{
		1: {TripCount: 10, AvgFare: 20.0, LastSeen: time.Now()},
		2: {TripCount: 5, AvgFare: 30.0, LastSeen: time.Now().Add(-10 * time.Minute)},
	}

	alerts := d.Detect(snapshots)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Zone != 2 {
		t.Errorf("expected alert for zone 2, got zone %d", alerts[0].Zone)
	}
	if alerts[0].Type != "DeadZone" {
		t.Errorf("expected type DeadZone, got %s", alerts[0].Type)
	}
}

func TestDeadZoneNoAlertForActiveZones(t *testing.T) {
	d := &DeadZoneDetector{InactiveThreshold: 5 * time.Minute}
	snapshots := map[int]aggregator.ZoneSnapshot{
		1: {TripCount: 10, AvgFare: 20.0, LastSeen: time.Now()},
		2: {TripCount: 5, AvgFare: 30.0, LastSeen: time.Now().Add(-2 * time.Minute)},
	}

	alerts := d.Detect(snapshots)

	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestMultipleDetectorsCompose(t *testing.T) {
	detectors := []Detector{
		&FareSpikeDetector{Threshold: 50.0},
		&DeadZoneDetector{InactiveThreshold: 5 * time.Minute},
	}
	snapshots := map[int]aggregator.ZoneSnapshot{
		1: {TripCount: 10, AvgFare: 80.0, LastSeen: time.Now().Add(-10 * time.Minute)},
	}

	var allAlerts []Alert
	for _, d := range detectors {
		allAlerts = append(allAlerts, d.Detect(snapshots)...)
	}

	if len(allAlerts) != 2 {
		t.Errorf("expected 2 alerts (fare spike + dead zone), got %d", len(allAlerts))
	}
}
