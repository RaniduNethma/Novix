package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/novix/services/processing-service/pkg/logger"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	return &Producer{writer: writer}
}

func (p *Producer) Publish(ctx context.Context,
	topic string, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	logger.Info("Published Kafka message",
		zap.String("topic", topic),
		zap.String("key", key),
	)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}