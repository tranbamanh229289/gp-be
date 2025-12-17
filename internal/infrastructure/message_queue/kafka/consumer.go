package kafka

import (
	"be/config"
	"be/pkg/logger"
	"context"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Handler func(ctx context.Context, msg *Message) error

type Consumer struct {
	config   *config.Config
	logger   *logger.ZapLogger
	consumer *kafka.Consumer
	handlers map[string]Handler
	mu       sync.RWMutex
	wg       sync.WaitGroup
}

func NewConsumer(config *config.Config, logger *logger.ZapLogger, groupID string) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       config.Kafka.Consumer.BootstrapServers,
		"group.id":                groupID,
		"auto.offset.reset":       config.Kafka.Consumer.AutoOffsetReset,
		"session.timeout.ms":      config.Kafka.Consumer.SessionTimeoutMs,
		"enable.auto.commit":      config.Kafka.Consumer.EnableAutoCommit,
		"auto.commit.interval.ms": config.Kafka.Consumer.AutoCommitIntervalMs,
	})
	return &Consumer{
		config:   config,
		logger:   logger,
		consumer: consumer,
		handlers: make(map[string]Handler),
	}, err
}

func (c *Consumer) Register(topic string, handler Handler) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[topic] = handler
	return nil
}

func (c *Consumer) ReceiveMessage(ctx context.Context) error {
	topics := make([]string, 0, len(c.handlers))
	for k, _ := range c.handlers {
		topics = append(topics, k)
	}
	err := c.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %s", err)
	}
	c.wg.Add(1)
	go c.consume(ctx)
	return nil
}

func (c *Consumer) consume(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				c.handleMessage(ctx, e)
			case kafka.Error:
				c.logger.Error(fmt.Sprintf("Kafka error: %s", e))
				if e.Code() == kafka.ErrAllBrokersDown {
					return
				}
			}
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg *kafka.Message) {
	topic := *msg.TopicPartition.Topic
	c.mu.RLock()
	defer c.mu.RUnlock()
	handler, exists := c.handlers[topic]

	consumerMsg := &Message{
		Topic:     topic,
		Key:       string(msg.Key),
		Value:     msg.Value,
		Headers:   make(map[string]string),
		Partition: msg.TopicPartition.Partition,
	}

	for _, h := range msg.Headers {
		consumerMsg.Headers[h.Key] = string(h.Value)
	}

	if !exists {
		c.logger.Error(fmt.Sprintf("No handler for topic:%s", topic))
		return
	}

	err := handler(ctx, consumerMsg)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Handler error for topic:%s", topic))
		return
	}

	if !c.config.Kafka.Consumer.EnableAutoCommit {
		_, err := c.consumer.CommitMessage(msg)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Commit Failed: %s", err))
			return
		}
	}
}

func (c *Consumer) Commit() error {
	_, err := c.consumer.Commit()
	return err
}

func (c *Consumer) Close(ctx context.Context) error {
	c.wg.Wait()
	return c.consumer.Close()
}
