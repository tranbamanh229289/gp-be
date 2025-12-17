package rabbitmq

import (
	"be/config"
	"be/pkg/logger"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitQueue struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	logger     *logger.ZapLogger
}

func NewQueue(cfg *config.Config, logger *logger.ZapLogger) (*RabbitQueue, error) {
	// connect
	dsn := config.GetRabbitMQDSN(cfg)
	conn, err := amqp091.Dial(dsn)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to RabbitMQ: %s", err))
		return nil, err
	}

	// open channel
	ch, err := conn.Channel()

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open channel RabbitMQ %s", err))
		conn.Close()
		return nil, err
	}

	// set qos
	if err := ch.Qos(cfg.RabbitMQ.PrefetchCount, 0, false); err != nil {
		logger.Error(fmt.Sprintf("Failed to set qos: %s", err))
		return nil, err
	}
	logger.Info("Successfully connected to RabbitMQ")
	return &RabbitQueue{connection: conn, channel: ch, logger: logger}, nil

}

func (rq *RabbitQueue) DeclareExchange(name, kind string, durable, autoDelete bool) error {
	err := rq.channel.ExchangeDeclare(name, kind, durable, autoDelete, false, false, nil)
	if err != nil {
		rq.logger.Error(fmt.Sprintf("Failed to declare exchange: %s", err))
		return err
	}
	return nil
}

func (rq *RabbitQueue) DeclareQueue(name string, durable, autoDelete, exclusive bool) (*amqp091.Queue, error) {
	queue, err := rq.channel.QueueDeclare(name, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		rq.logger.Error(fmt.Sprintf("Failed to declare queue: %s", err))
		return &queue, err
	}
	return &queue, nil
}

func (rq *RabbitQueue) BindQueue(queue, routingKey, exchange string) error {
	err := rq.channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		rq.logger.Error(fmt.Sprintf("Failed to bind queue: %s", err))
		return err
	}
	return nil
}

func (rq *RabbitQueue) UnbindQueue(queue, routingKey, exchange string) error {
	err := rq.channel.QueueUnbind(queue, routingKey, exchange, nil)
	if err != nil {
		rq.logger.Error(fmt.Sprintf("Failed to unbind queue: %s", err))
		return err
	}
	return nil
}

func (rq *RabbitQueue) Close() error {
	if rq.channel != nil {
		if err := rq.channel.Close(); err != nil {
			rq.logger.Error(fmt.Sprintf("Failed to close channel: %s", err))
		}
	}
	if rq.connection != nil {
		if err := rq.connection.Close(); err != nil {
			rq.logger.Error(fmt.Sprintf("Failed to close connection: %s", err))
		}
	}
	return nil
}
