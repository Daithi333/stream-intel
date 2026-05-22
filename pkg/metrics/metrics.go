package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	eventsTotal prometheus.Counter
	errorsTotal prometheus.Counter
	consumerLag prometheus.Gauge
}

func New() *Metrics {
	m := &Metrics{
		eventsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "stream_intel_events_total",
			Help: "Total events consumed from Kafka",
		}),
		errorsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "stream_intel_consumer_errors_total",
			Help: "Total consumer errors",
		}),
		consumerLag: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "stream_intel_consumer_lag",
			Help: "Consumer group lag",
		}),
	}
	prometheus.MustRegister(m.eventsTotal, m.errorsTotal, m.consumerLag)
	return m
}

func (m *Metrics) RecordEvent() {
	m.eventsTotal.Inc()
}

func (m *Metrics) RecordError() {
	m.errorsTotal.Inc()
}

func (m *Metrics) Run(ctx context.Context, port string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}
