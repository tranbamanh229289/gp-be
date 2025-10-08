package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	client *RabbitMQClient
}

func NewProducer(client *RabbitMQClient)(*Producer, error) {
	return &Producer{client: client}, nil
}

func(p *Producer) Publish(ctx context.Context, exchange string, message Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Failed to marshal message: %w", err)
	}
	
	err = p.client.Channel.Publish(exchange, message.RoutingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: body,
		DeliveryMode: amqp091.Persistent,
		Timestamp: time.Now(),
		Headers: message.Headers,
		Priority: message.Priority,
	})

	if err != nil {
		log.Fatal("Failed to publish: %w", err)
	}

	return nil
}

func (p *Producer) PublishWithAck(ctx context.Context, exchange, routingKey string, message Message) error {
	if err := p.client.Channel.Confirm(false); err != nil {
		log.Fatal("Failed to put channel in confirm mode: %w", err)
	}

	confirms := p.client.Channel.NotifyPublish(make(chan amqp091.Confirmation, 1))

	body, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Failed to marshal message: %w", err)
	}
	
		err = p.client.Channel.Publish(exchange, message.RoutingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: body,
		DeliveryMode: amqp091.Persistent,
		Timestamp: time.Now(),
		Headers: message.Headers,
		Priority: message.Priority,
	})

	if err != nil {
		log.Fatal("Failed to publish: %w", err)
	}

	confirmed := <- confirms
	if !confirmed.Ack {
		log.Fatal("Failed to receive confirm")
	}

	return nil 
}