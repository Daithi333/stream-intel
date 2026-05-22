# stream-intel

A realtime stream intelligence service written in Go, designed as a companion to [data-streaming-platform](https://github.com/dmcelhill/data-streaming-platform).

Consumes the same Kafka topic (`taxi_trips`) independently from Spark, providing low-latency operational analytics alongside the batch/lakehouse medallion architecture.


---

## Purpose

This project is a hands-on exploration of Go systems programming concepts through a realistic streaming workload:

- Goroutines and worker pools (concurrent partition consumers)
- Channels and select (internal event pipeline, backpressure)
- Context propagation and cancellation (graceful shutdown)
- Pointers and memory management (shared aggregation state)
- Interfaces (pluggable detectors, sinks)
- Synchronisation primitives (mutexes, atomics, errgroup)


---

## Architecture

```
Kafka (taxi_trips topic)
    |
    v
[Partition Consumers]  -- 1 goroutine per partition
    |
    | (bounded channels)
    v
[Aggregation Workers]  -- sliding window stats by zone
    |
    +---> [Anomaly Detector]  -- fare spikes, dead zones
    |
    +---> [Prometheus Exporter]  -- /metrics endpoint
    |
    +---> [WebSocket Server]  -- live dashboard feed
```

The service runs as an independent Kafka consumer group. It does not interfere with Spark's Bronze consumer — both read from the same topic via separate group IDs.


---

## Planned Capabilities

### Phase 1 — Consumer and Metrics
- Concurrent Kafka partition consumers (consumer group rebalancing)
- Internal channel-based event pipeline
- Prometheus metrics (events/sec, consumer lag, partition assignment)
- Graceful shutdown via context cancellation and signal handling
- Structured logging
- Configuration via environment variables

### Phase 2 — Aggregations and Anomaly Detection
- Sliding window aggregations (trips/min, avg fare by pickup zone)
- Pluggable anomaly detection interface (fare spikes, inactive zones)
- In-memory state with mutex-protected access
- Alerting via configurable sinks

### Phase 3 — WebSocket Dashboard and Replay
- WebSocket endpoint for live aggregation feed
- Fan-out pattern to multiple connected clients
- Replay mode (consume from earliest offset on demand)
- Backpressure strategies (drop, buffer, slow consumer)


---

## Project Structure

```
stream-intel/
  cmd/
    stream-intel/       # Application entrypoint
      main.go
  internal/
    config/             # Environment-based configuration
    consumer/           # Kafka consumer group, partition workers
    pipeline/           # Channel-based event routing
    aggregator/         # Windowed statistics
    detector/           # Anomaly detection interface + implementations
    metrics/            # Prometheus instrumentation
    websocket/          # WebSocket server and fan-out
  pkg/
    model/              # Shared event types
  deployments/
    docker-compose.yml  # Local dev (points to shared Kafka)
  Makefile
  go.mod
  go.sum
  .gitignore
  README.md
```

Follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) conventions:
- `cmd/` for application binaries
- `internal/` for private application code (not importable by other modules)
- `pkg/` for code that could be imported by external projects


---

## Prerequisites

- Go 1.22+
- Access to the Kafka broker from `data-streaming-platform` (default: `localhost:9092`)
- Topic `taxi_trips` with events being produced


---

## Quick Start

```bash
# Start the shared Kafka broker (from data-streaming-platform)
cd ../data-streaming-platform && make up && make topic-create

# Run the service
go run ./cmd/stream-intel

# Produce events (from data-streaming-platform)
cd ../data-streaming-platform && make produce
```


---

## Configuration

All configuration via environment variables with sensible defaults:

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFKA_BROKERS` | `localhost:9092` | Comma-separated broker addresses |
| `KAFKA_TOPIC` | `taxi_trips` | Topic to consume |
| `KAFKA_GROUP_ID` | `stream-intel` | Consumer group ID |
| `METRICS_PORT` | `9090` | Prometheus metrics HTTP port |
| `WS_PORT` | `8080` | WebSocket server port |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |


---

## Key Dependencies

| Library | Purpose |
|---------|---------|
| [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) | Kafka consumer (librdkafka-based, production grade) |
| [prometheus/client_golang](https://github.com/prometheus/client_golang) | Metrics exposition |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket connections |
| [slog](https://pkg.go.dev/log/slog) | Structured logging (stdlib, Go 1.21+) |


---

## Relationship to data-streaming-platform

```
[Event Producer] --> [Kafka: taxi_trips]
                          |
              +-----------+-----------+
              |                       |
     [Spark Bronze/Silver/Gold]   [stream-intel]
     (batch lakehouse, Delta)     (realtime ops)
```

Both consumers operate independently. Spark owns the durable analytical path (medallion architecture, Delta Lake). This service owns the operational/realtime path (dashboards, alerts, metrics).


---

## Development Roadmap

- [ ] Phase 1: Consumer + Prometheus metrics
- [ ] Phase 2: Windowed aggregations + anomaly detection
- [ ] Phase 3: WebSocket live feed + replay support
- [ ] CI pipeline (lint, test, build)
- [ ] Docker image for deployment alongside the platform
