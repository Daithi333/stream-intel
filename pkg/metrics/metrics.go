package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	eventsTotal     prometheus.Counter
	errorsTotal     prometheus.Counter
	consumerLag     prometheus.Gauge
	pipelineBacklog prometheus.Gauge
	droppedMessages prometheus.Counter
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
		pipelineBacklog: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "stream_intel_pipeline_backlog",
			Help: "Number of events buffered in the pipeline channel",
		}),
		droppedMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "stream_intel_ws_dropped_messages_total",
			Help: "Total messages dropped due to slow WebSocket clients",
		}),
	}
	prometheus.MustRegister(m.eventsTotal, m.errorsTotal, m.consumerLag, m.pipelineBacklog, m.droppedMessages)
	return m
}

func (m *Metrics) RecordEvent() {
	m.eventsTotal.Inc()
}

func (m *Metrics) RecordError() {
	m.errorsTotal.Inc()
}

func (m *Metrics) SetPipelineBacklog(size int) {
	m.pipelineBacklog.Set(float64(size))
}

func (m *Metrics) RecordDroppedMessage() {
	m.droppedMessages.Inc()
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
		_ = server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}
