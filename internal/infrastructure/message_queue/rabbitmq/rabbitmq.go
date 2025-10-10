package rabbitmq

import (
	"be/config"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitQueue struct{
	Connection *amqp091.Connection
	Channel *amqp091.Channel
}


func NewQueue(cfg *config.Config) (*RabbitQueue, error){
	// connect
	dsn := config.GetRabbitMQDSN(cfg)
	conn, err := amqp091.Dial(dsn)

	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: %w", err)
	}

	// open channel
	ch, err := conn.Channel()

	if err != nil {
		conn.Close()
		log.Fatal("Failed to open channel RabbitMQ")
	}

	// set qos
	if err := ch.Qos(cfg.RabbitMQ.PrefetchCount, 0, false); err != nil {
		log.Fatal("Failed to set qos %w", err)
	}

	return &RabbitQueue{Connection: conn, Channel: ch}, nil 

}

func (rq *RabbitQueue) DeclareExchange(name, kind string, durable, autoDelete bool) error {
	err := rq.Channel.ExchangeDeclare(name, kind, durable, autoDelete, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange: %w", err)
	}
	return nil
}	


func (rq *RabbitQueue) DeclareQueue(name string, durable, autoDelete, exclusive bool) (*amqp091.Queue, error) {
	queue, err := rq.Channel.QueueDeclare(name, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue: %w", err)
	}
	return &queue, nil 
}

func (rq *RabbitQueue) BindQueue(queue, routingKey, exchange string) error {
	err := rq.Channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		log.Fatal("Failed to bind queue: %w", err)
	}
	return nil
}

func (rq *RabbitQueue) UnbindQueue(queue, routingKey, exchange string) error {
	err := rq.Channel.QueueUnbind(queue, routingKey, exchange, nil)
	if err != nil {
		log.Fatal("Failed to unbind queue: %w", err)
	}
	return nil
}

func (rq *RabbitQueue) Close()error {
	if rq.Channel != nil {
		if err := rq.Channel.Close(); err != nil {
			log.Fatal("Failed to close channel :%w", err)
		}
	}
	if rq.Connection != nil {
		if err := rq.Connection.Close(); err != nil {
			log.Fatal("Failed to close connection : %w", err)
		}
	}
	return nil
}
