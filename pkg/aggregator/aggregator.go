package aggregator

import (
	"sync"
	"time"

	"github.com/dmcelhill/stream-intel/pkg/model"
)

type ZoneStats struct {
	TripCount int
	FareSum   float64
	FirstSeen time.Time
	LastSeen  time.Time
}

type ZoneSnapshot struct {
	TripCount int
	AvgFare   float64
	LastSeen  time.Time
}

type Aggregator struct {
	mu    sync.RWMutex
	zones map[int]*ZoneStats
}

func New() *Aggregator {
	return &Aggregator{
		zones: make(map[int]*ZoneStats),
	}
}

func (a *Aggregator) Record(trip model.TaxiTrip) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if stats, exists := a.zones[trip.PuLocationId]; exists {
		stats.TripCount++
		stats.FareSum += trip.FareAmount
		stats.LastSeen = time.Now()
	} else {
		a.zones[trip.PuLocationId] = &ZoneStats{
			TripCount: 1,
			FareSum:   trip.FareAmount,
			FirstSeen: time.Now(),
			LastSeen:  time.Now(),
		}
	}
}

func (a *Aggregator) Snapshot() map[int]ZoneSnapshot {
	a.mu.RLock()
	defer a.mu.RUnlock()
	snapshot := make(map[int]ZoneSnapshot)
	for zone, stats := range a.zones {
		snapshot[zone] = ZoneSnapshot{
			TripCount: stats.TripCount,
			AvgFare:   stats.FareSum / float64(stats.TripCount),
			LastSeen:  stats.LastSeen,
		}
	}
	return snapshot
}
