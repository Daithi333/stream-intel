package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	KafkaBrokers       []string
	KafkaTopic         string
	KafkaGroupID       string
	MetricsPort        string
	WSPort             string
	LogLevel           string
	PipelineBufferSize int
	DetectorInterval   int
	FareSpikeThreshold float64
	DeadZoneThreshold  int
}

func Load() Config {
	return Config{
		KafkaBrokers:       strings.Split(getEnv("KAFKA_BROKERS", "localhost:9094"), ","),
		KafkaTopic:         getEnv("KAFKA_TOPIC", "taxi_trips"),
		KafkaGroupID:       getEnv("KAFKA_GROUP_ID", "stream-intel"),
		MetricsPort:        getEnv("METRICS_PORT", "9090"),
		WSPort:             getEnv("WS_PORT", "8080"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		PipelineBufferSize: getEnvInt("PIPELINE_BUFFER_SIZE", 100),
		DetectorInterval:   getEnvInt("DETECTOR_INTERVAL_SECS", 10),
		FareSpikeThreshold: getEnvFloat("FARE_SPIKE_THRESHOLD", 50.0),
		DeadZoneThreshold:  getEnvInt("DEAD_ZONE_THRESHOLD_SECS", 300),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
