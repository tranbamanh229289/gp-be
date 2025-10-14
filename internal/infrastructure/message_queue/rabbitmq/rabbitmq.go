package rabbitmq

import (
	"be/config"
	"be/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitQueue struct{
	connection *amqp091.Connection
	channel *amqp091.Channel
	logger *logger.ZapLogger
}


func NewQueue(cfg *config.Config, logger *logger.ZapLogger) (*RabbitQueue, error){
	// connect
	dsn := config.GetRabbitMQDSN(cfg)
	conn, err := amqp091.Dial(dsn)

	if err != nil {
		logger.Error("Failed to connect to RabbitMQ: ", zap.Error(err), zap.String("addresses", dsn))
		return nil, err
	}

	// open channel
	ch, err := conn.Channel()

	if err != nil {
		logger.Error("Failed to open channel RabbitMQ", zap.Error(err))
		conn.Close()
		return nil, err
	}

	// set qos
	if err := ch.Qos(cfg.RabbitMQ.PrefetchCount, 0, false); err != nil {
		logger.Error("Failed to set qos", zap.Error(err))
		return nil, err
	}

	return &RabbitQueue{connection: conn, channel: ch, logger: logger}, nil 

}

func (rq *RabbitQueue) DeclareExchange(name, kind string, durable, autoDelete bool) error {
	err := rq.channel.ExchangeDeclare(name, kind, durable, autoDelete, false, false, nil)
	if err != nil {
		rq.logger.Error("Failed to declare exchange:", zap.Error(err))
		return err
	}
	return nil
}	


func (rq *RabbitQueue) DeclareQueue(name string, durable, autoDelete, exclusive bool) (*amqp091.Queue, error) {
	queue, err := rq.channel.QueueDeclare(name, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		rq.logger.Error("Failed to declare queue:", zap.Error(err))
		return &queue, err
	}
	return &queue, nil 
}

func (rq *RabbitQueue) BindQueue(queue, routingKey, exchange string) error {
	err := rq.channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		rq.logger.Error("Failed to bind queue:", zap.Error(err))
		return err
	}
	return nil
}

func (rq *RabbitQueue) UnbindQueue(queue, routingKey, exchange string) error {
	err := rq.channel.QueueUnbind(queue, routingKey, exchange, nil)
	if err != nil {
		rq.logger.Error("Failed to unbind queue: ", zap.Error(err))
		return err
	}
	return nil
}

func (rq *RabbitQueue) Close()error {
	if rq.channel != nil {
		if err := rq.channel.Close(); err != nil {
			rq.logger.Error("Failed to close channel: ", zap.Error(err))
		}
	}
	if rq.connection != nil {
		if err := rq.connection.Close(); err != nil {
			rq.logger.Error("Failed to close connection:", zap.Error(err))
		}
	}
	return nil
}
