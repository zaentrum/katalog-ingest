// Package kafka wraps a sarama sync producer + consumer group used by
// katalog-ingest. The scheduler emits stube.processing.task.* messages; the
// consumer side watches stube.processing.result.* for downstream completion
// events so the scheduler can mark items done.
package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IBM/sarama"
)

// Producer wraps a sarama.SyncProducer. A nil *Producer is a permitted
// scaffold-mode value; Publish becomes a no-op.
type Producer struct {
	p sarama.SyncProducer
}

// NewProducer connects to the given comma-separated broker list.
func NewProducer(brokers string) (*Producer, error) {
	if brokers == "" {
		return nil, errors.New("kafka brokers empty")
	}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	// TODO: parse brokers as []string, pull tls/sasl from env.
	_ = cfg
	return &Producer{}, nil
}

// Publish writes a single record. Scaffold: logs and returns nil.
func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	if p == nil || p.p == nil {
		slog.Debug("kafka publish (scaffold no-op)", "topic", topic, "bytes", len(value))
		return nil
	}
	_, _, err := p.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	return err
}

// Close releases the producer.
func (p *Producer) Close() error {
	if p == nil || p.p == nil {
		return nil
	}
	return p.p.Close()
}
