package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/novix/services/processing-service/pkg/logger"
	"go.uber.org/zap"
)

type VideoUploadedEvent struct {
	VideoID string `json:"videoId"`
	UserID  string `json:"userId"`
	RawPath string `json:"rawPath"`
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{reader: reader}
}

// Consume reads messages and calls the handler function
func (c *Consumer) Consume(ctx context.Context,
	handler func(event VideoUploadedEvent) error) error {
	logger.Info("Starting Kafka consumer")
	for {
		select {
		case <-ctx.Done():
			logger.Info("Kafka consumer shutting down")
			return nil
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				logger.Error("Failed to read Kafka message",
					zap.Error(err))
				continue
			}
			var event VideoUploadedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Error("Failed to unmarshal event",
					zap.Error(err),
					zap.String("raw", string(msg.Value)),
				)
				continue
			}
			logger.Info("Received video uploaded event",
				zap.String("videoID", event.VideoID),
				zap.String("userID", event.UserID),
			)
			if err := handler(event); err != nil {
				logger.Error("Failed to handle event",
					zap.String("videoID", event.VideoID),
					zap.Error(err),
				)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return fmt.Errorf("consumer closed: %w", c.reader.Close())
}