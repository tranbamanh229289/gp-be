package rabbitmq

import (
	"be/pkg/logger"
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Producer struct {
	queue *RabbitQueue
	logger *logger.ZapLogger
}

func NewProducer(queue *RabbitQueue, logger *logger.ZapLogger)(*Producer, error) {
	return &Producer{queue: queue, logger: logger}, nil
}

func(p *Producer) Publish(ctx context.Context, exchange string, message Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("Failed to marshal message:", zap.Error(err))
		return err
	}
	
	err = p.queue.channel.Publish(exchange, message.RoutingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: body,
		DeliveryMode: amqp091.Persistent,
		Timestamp: time.Now(),
		Headers: message.Headers,
		Priority: message.Priority,
	})

	if err != nil {
		p.logger.Error("Failed to publish:", zap.Error(err))
		return err
	}

	return nil
}

func (p *Producer) PublishWithAck(ctx context.Context, exchange, routingKey string, message Message) error {
	if err := p.queue.channel.Confirm(false); err != nil {
		p.logger.Error("Failed to put channel in confirm mode:", zap.Error(err))
		return err
	}

	confirms := p.queue.channel.NotifyPublish(make(chan amqp091.Confirmation, 1))

	body, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("Failed to marshal message:", zap.Error(err))
		return err
	}
	
		err = p.queue.channel.Publish(exchange, message.RoutingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: body,
		DeliveryMode: amqp091.Persistent,
		Timestamp: time.Now(),
		Headers: message.Headers,
		Priority: message.Priority,
	})

	if err != nil {
		p.logger.Error("Failed to publish:", zap.Error(err))
		return err
	}

	confirmed := <- confirms
	if !confirmed.Ack {
		p.logger.Error("Failed to receive confirm")
	}

	return nil 
}