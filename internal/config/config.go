package config

import (
	"os"
	"strings"
)

type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
	MetricsPort  string
	WSPort       string
	LogLevel     string
}

func Load() Config {
	return Config{
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "taxi_trips"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "stream-intel"),
		MetricsPort:  getEnv("METRICS_PORT", "9090"),
		WSPort:       getEnv("WS_PORT", "8080"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
