package aggregator

import (
	"testing"
	"time"

	"github.com/dmcelhill/stream-intel/pkg/model"
)

func TestRecordNewZone(t *testing.T) {
	agg := New()

	agg.Record(model.TaxiTrip{PuLocationId: 42, FareAmount: 20.0})

	snap := agg.Snapshot()
	zone, exists := snap[42]
	if !exists {
		t.Fatal("expected zone 42 to exist in snapshot")
	}
	if zone.TripCount != 1 {
		t.Errorf("expected TripCount 1, got %d", zone.TripCount)
	}
	if zone.AvgFare != 20.0 {
		t.Errorf("expected AvgFare 20.0, got %f", zone.AvgFare)
	}
}

func TestRecordMultipleTrips(t *testing.T) {
	agg := New()

	agg.Record(model.TaxiTrip{PuLocationId: 7, FareAmount: 10.0})
	agg.Record(model.TaxiTrip{PuLocationId: 7, FareAmount: 30.0})
	agg.Record(model.TaxiTrip{PuLocationId: 7, FareAmount: 20.0})

	snap := agg.Snapshot()
	zone := snap[7]
	if zone.TripCount != 3 {
		t.Errorf("expected TripCount 3, got %d", zone.TripCount)
	}
	if zone.AvgFare != 20.0 {
		t.Errorf("expected AvgFare 20.0, got %f", zone.AvgFare)
	}
}

func TestRecordMultipleZones(t *testing.T) {
	agg := New()

	agg.Record(model.TaxiTrip{PuLocationId: 1, FareAmount: 15.0})
	agg.Record(model.TaxiTrip{PuLocationId: 2, FareAmount: 25.0})

	snap := agg.Snapshot()
	if len(snap) != 2 {
		t.Errorf("expected 2 zones, got %d", len(snap))
	}
	if snap[1].AvgFare != 15.0 {
		t.Errorf("expected zone 1 AvgFare 15.0, got %f", snap[1].AvgFare)
	}
	if snap[2].AvgFare != 25.0 {
		t.Errorf("expected zone 2 AvgFare 25.0, got %f", snap[2].AvgFare)
	}
}

func TestSnapshotIsIndependent(t *testing.T) {
	agg := New()
	agg.Record(model.TaxiTrip{PuLocationId: 1, FareAmount: 10.0})

	snap := agg.Snapshot()

	agg.Record(model.TaxiTrip{PuLocationId: 1, FareAmount: 90.0})

	if snap[1].TripCount != 1 {
		t.Error("snapshot should not be affected by subsequent records")
	}
}

func TestLastSeenUpdates(t *testing.T) {
	agg := New()

	agg.Record(model.TaxiTrip{PuLocationId: 5, FareAmount: 10.0})
	time.Sleep(10 * time.Millisecond)
	agg.Record(model.TaxiTrip{PuLocationId: 5, FareAmount: 20.0})

	snap := agg.Snapshot()
	if time.Since(snap[5].LastSeen) > time.Second {
		t.Error("LastSeen should be recent")
	}
}
