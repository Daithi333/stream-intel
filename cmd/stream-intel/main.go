package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmcelhill/stream-intel/internal/config"
	"github.com/dmcelhill/stream-intel/pkg/aggregator"
	"github.com/dmcelhill/stream-intel/pkg/consumer"
	"github.com/dmcelhill/stream-intel/pkg/detector"
	"github.com/dmcelhill/stream-intel/pkg/metrics"
	"github.com/dmcelhill/stream-intel/pkg/pipeline"
	"github.com/dmcelhill/stream-intel/pkg/sink"
	ws "github.com/dmcelhill/stream-intel/pkg/websocket"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	pipe := pipeline.New(cfg.PipelineBufferSize)
	c, err := consumer.New(cfg, pipe, logger)
	if err != nil {
		logger.Error("Failed to initialise Message Consumer")
		os.Exit(1)
	}
	m := metrics.New()
	agg := aggregator.New()
	detectors := []detector.Detector{
		&detector.FareSpikeDetector{Threshold: cfg.FareSpikeThreshold},
		&detector.DeadZoneDetector{InactiveThreshold: time.Duration(cfg.DeadZoneThreshold) * time.Second},
	}

	hub := ws.NewHub()
	hub.OnDrop = m.RecordDroppedMessage
	wsSrv := ws.NewServer(hub, logger)
	wsSrv.SetReplayFunc(c.Replay)

	sinks := []sink.Sink{
		&sink.LogSink{Logger: logger},
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	logger.Info("stream-intel started",
		"brokers", cfg.KafkaBrokers,
		"topic", cfg.KafkaTopic,
		"metrics_port", cfg.MetricsPort,
		"ws_port", cfg.WSPort,
	)

	g.Go(func() error { return c.Run(ctx) })
	g.Go(func() error { return m.Run(ctx, cfg.MetricsPort) })
	g.Go(func() error { return wsSrv.Run(ctx, cfg.WSPort) })
	g.Go(processEvents(pipe, agg, m))
	g.Go(runDetectors(ctx, agg, detectors, sinks, time.Duration(cfg.DetectorInterval)*time.Second))
	g.Go(broadcastSnapshots(ctx, agg, hub, time.Duration(cfg.DetectorInterval)*time.Second))
	g.Go(shutdownOnCancel(ctx, pipe))

	if err := g.Wait(); err != nil {
		logger.Error("Service exiting", "error", err)
	}

	c.Close()
}

func processEvents(pipe *pipeline.Pipeline, agg *aggregator.Aggregator, m *metrics.Metrics) func() error {
	return func() error {
		for trip := range pipe.Receive() {
			agg.Record(trip)
			m.RecordEvent()
			m.SetPipelineBacklog(pipe.Len())
		}
		return nil
	}
}

func runDetectors(ctx context.Context, agg *aggregator.Aggregator, detectors []detector.Detector, sinks []sink.Sink, interval time.Duration) func() error {
	return func() error {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				snap := agg.Snapshot()
				for _, d := range detectors {
					alerts := d.Detect(snap)
					for _, alert := range alerts {
						for _, s := range sinks {
							_ = s.Send(alert)
						}
					}
				}
			}
		}
	}
}

func broadcastSnapshots(ctx context.Context, agg *aggregator.Aggregator, hub *ws.Hub, interval time.Duration) func() error {
	return func() error {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				snap := agg.Snapshot()
				data, err := json.Marshal(snap)
				if err != nil {
					continue
				}
				hub.Broadcast(data)
			}
		}
	}
}

func shutdownOnCancel(ctx context.Context, pipe *pipeline.Pipeline) func() error {
	return func() error {
		<-ctx.Done()
		pipe.Close()
		return nil
	}
}
