package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dmcelhill/stream-intel/internal/config"
	"github.com/dmcelhill/stream-intel/pkg/model"
	"github.com/dmcelhill/stream-intel/pkg/pipeline"
)

type Consumer struct {
	client   *kafka.Consumer
	pipeline *pipeline.Pipeline
	logger   *slog.Logger
}

func New(cfg config.Config, p *pipeline.Pipeline, logger *slog.Logger) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.KafkaBrokers, ","),
		"group.id":          cfg.KafkaGroupID,
	})
	if err != nil {
		return nil, err
	}
	err = c.SubscribeTopics([]string{cfg.KafkaTopic}, nil)
	if err != nil {
		c.Close()
		return nil, err
	}

	_, err = c.GetMetadata(nil, true, 5000)
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("broker unreachable: %w", err)
	}

	return &Consumer{
		client:   c,
		pipeline: p,
		logger:   logger,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	consecutiveErrors := 0
	maxErrors := 10

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			message, err := c.client.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
					continue
				}
				consecutiveErrors++
				c.logger.Error("Failed to read message", "error", err, "attempt", consecutiveErrors)
				if consecutiveErrors >= maxErrors {
					return fmt.Errorf("consumer error after %d attempts: %w", consecutiveErrors, err)
				}
				time.Sleep(time.Duration(consecutiveErrors) * time.Second)
				continue
			}
			consecutiveErrors = 0

			var trip model.TaxiTrip
			if err := json.Unmarshal(message.Value, &trip); err != nil {
				c.logger.Error("Failed to unmarshal message", "error", err)
				continue
			}
			if err := c.pipeline.Send(ctx, trip); err != nil {
				c.logger.Error("Failed to send to pipeline", "error", err)
			}
		}
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}

func (c *Consumer) Replay() error {
	assignment, err := c.client.Assignment()
	if err != nil {
		return fmt.Errorf("failed to get partition assignment: %w", err)
	}
	for i := range assignment {
		assignment[i].Offset = kafka.OffsetBeginning
	}
	err = c.client.Assign(assignment)
	if err != nil {
		return fmt.Errorf("failed to seek to beginning: %w", err)
	}
	c.logger.Info("Replay initiated, seeking to beginning of all partitions")
	return nil
}
