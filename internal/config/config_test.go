package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()

	cfg := Load()

	if len(cfg.KafkaBrokers) != 1 || cfg.KafkaBrokers[0] != "localhost:9092" {
		t.Errorf("expected KafkaBrokers [localhost:9092], got %v", cfg.KafkaBrokers)
	}
	if cfg.KafkaTopic != "taxi_trips" {
		t.Errorf("expected KafkaTopic taxi_trips, got %s", cfg.KafkaTopic)
	}
	if cfg.KafkaGroupID != "stream-intel" {
		t.Errorf("expected KafkaGroupID stream-intel, got %s", cfg.KafkaGroupID)
	}
	if cfg.MetricsPort != "9090" {
		t.Errorf("expected MetricsPort 9090, got %s", cfg.MetricsPort)
	}
	if cfg.WSPort != "8080" {
		t.Errorf("expected WSPort 8080, got %s", cfg.WSPort)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected LogLevel info, got %s", cfg.LogLevel)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("KAFKA_BROKERS", "broker1:9092,broker2:9092")
	t.Setenv("KAFKA_TOPIC", "custom_topic")
	t.Setenv("KAFKA_GROUP_ID", "custom-group")
	t.Setenv("METRICS_PORT", "9191")
	t.Setenv("WS_PORT", "8181")
	t.Setenv("LOG_LEVEL", "debug")

	cfg := Load()

	if len(cfg.KafkaBrokers) != 2 || cfg.KafkaBrokers[0] != "broker1:9092" || cfg.KafkaBrokers[1] != "broker2:9092" {
		t.Errorf("expected KafkaBrokers [broker1:9092 broker2:9092], got %v", cfg.KafkaBrokers)
	}
	if cfg.KafkaTopic != "custom_topic" {
		t.Errorf("expected KafkaTopic custom_topic, got %s", cfg.KafkaTopic)
	}
	if cfg.KafkaGroupID != "custom-group" {
		t.Errorf("expected KafkaGroupID custom-group, got %s", cfg.KafkaGroupID)
	}
	if cfg.MetricsPort != "9191" {
		t.Errorf("expected MetricsPort 9191, got %s", cfg.MetricsPort)
	}
	if cfg.WSPort != "8181" {
		t.Errorf("expected WSPort 8181, got %s", cfg.WSPort)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel debug, got %s", cfg.LogLevel)
	}
}
