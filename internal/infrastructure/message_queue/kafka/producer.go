package kafka

import (
	"be/config"
	"be/pkg/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	config   *config.Config
	logger   *logger.ZapLogger
	producer *kafka.Producer
}

func NewProducer(config *config.Config, logger *logger.ZapLogger, clientID string) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  config.Kafka.Producer.BootstrapServers,
		"client.id":          clientID,
		"acks":               config.Kafka.Producer.Acks,
		"compression.type":   config.Kafka.Producer.CompressionType,
		"linger.ms":          config.Kafka.Producer.LingerMs,
		"batch.size":         config.Kafka.Producer.BatchSize,
		"enable.idempotence": config.Kafka.Producer.EnableIdempotence,
		"reties":             config.Kafka.Producer.Retries,
		"retry.backoff.ms":   100,
	})
	p := &Producer{
		config:   config,
		logger:   logger,
		producer: producer,
	}

	go p.handleDeliveryReports()
	return p, err
}

func (p *Producer) SendMessage(msg *Message) error {
	value, err := p.serialize(msg.Value)

	if err != nil {
		return fmt.Errorf("serialize error: %w", err)
	}

	partition := kafka.PartitionAny
	if msg.Partition > 0 {
		partition = msg.Partition
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &msg.Topic,
			Partition: partition,
		},
		Key:   []byte(msg.Key),
		Value: value,
	}

	if len(msg.Headers) > 0 {
		headers := make([]kafka.Header, 0, len(msg.Headers))
		for k, v := range msg.Headers {
			headers = append(kafkaMsg.Headers, kafka.Header{
				Key:   k,
				Value: []byte(v),
			})
		}
		kafkaMsg.Headers = headers
	}

	return p.producer.Produce(kafkaMsg, nil)
}

func (p *Producer) SendMessageSync(ctx context.Context, msg *Message) (*DeliveryReport, error) {
	deliveryChan := make(chan kafka.Event, 1)

	value, err := p.serialize(msg.Value)

	if err != nil {
		return nil, fmt.Errorf("serialize error: %w", err)
	}

	parition := kafka.PartitionAny
	if msg.Partition > 0 {
		parition = msg.Partition
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &msg.Topic,
			Partition: parition,
		},
		Key:   []byte(msg.Key),
		Value: value,
	}

	if len(msg.Headers) > 0 {
		kafkaMsg.Headers = make([]kafka.Header, 0, len(msg.Headers))
		for k, v := range msg.Headers {
			kafkaMsg.Headers = append(kafkaMsg.Headers, kafka.Header{
				Key:   k,
				Value: []byte(v),
			})
		}
	}
	err = p.producer.Produce(kafkaMsg, deliveryChan)

	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("send timeout after %s", ctx.Err())

	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		report := &DeliveryReport{
			Topic:     *m.TopicPartition.Topic,
			Partition: m.TopicPartition.Partition,
			Offset:    int64(m.TopicPartition.Offset),
			Key:       string(m.Key),
			Error:     m.TopicPartition.Error,
		}
		return report, m.TopicPartition.Error
	}
}

func (p *Producer) SendJSON(topic, key string, data interface{}) error {
	value, err := p.serialize(data)
	if err != nil {
		return fmt.Errorf("serialize error: %w", err)
	}
	return p.SendMessage(&Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Headers: map[string]string{
			"content-type": "application/json",
		},
	})
}

func (p *Producer) SendJSONSync(ctx context.Context, topic, key string, data interface{}) (*DeliveryReport, error) {
	value, err := p.serialize(data)
	if err != nil {
		return nil, fmt.Errorf("serialize error: %w", err)
	}
	return p.SendMessageSync(ctx, &Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Headers: map[string]string{
			"content-type": "application/json",
		},
	})
}

func (p *Producer) SendString(topic, key, data string) error {
	value, err := p.serialize(data)
	if err != nil {
		return fmt.Errorf("serialize error: %w", err)
	}
	return p.SendMessage(&Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

func (p *Producer) SendStringSync(ctx context.Context, topic, key, data string) (*DeliveryReport, error) {
	value, err := p.serialize(data)
	if err != nil {
		return nil, fmt.Errorf("serialize error: %w", err)
	}
	return p.SendMessageSync(ctx, &Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

func (p *Producer) Flush(timeout int) int {
	return p.producer.Flush(timeout)
}

func (p *Producer) handleDeliveryReports() {
	for e := range p.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				p.logger.Error(fmt.Sprintf("Delivery failed:%s", ev.TopicPartition.Error))
			}
		}
	}
}

func (p *Producer) serialize(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return json.Marshal(v)
	}
}
