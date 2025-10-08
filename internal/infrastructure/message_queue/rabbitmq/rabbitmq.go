package rabbitmq

import (
	"graduate-project/config"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct{
	Connection *amqp091.Connection
	Channel *amqp091.Channel
}


func NewClient(cfg *config.Config) (*RabbitMQClient, error){
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

	return &RabbitMQClient{Connection: conn, Channel: ch}, nil 

}

func (rc *RabbitMQClient) DeclareExchange(name, kind string, durable, autoDelete bool) error {
	err := rc.Channel.ExchangeDeclare(name, kind, durable, autoDelete, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange: %w", err)
	}
	return nil
}	


func (rc *RabbitMQClient) DeclareQueue(name string, durable, autoDelete, exclusive bool) (*amqp091.Queue, error) {
	queue, err := rc.Channel.QueueDeclare(name, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue: %w", err)
	}
	return &queue, nil 
}

func (rc *RabbitMQClient) BindQueue(queue, routingKey, exchange string) error {
	err := rc.Channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		log.Fatal("Failed to bind queue: %w", err)
	}
	return nil
}

func (rc *RabbitMQClient) UnbindQueue(queue, routingKey, exchange string) error {
	err := rc.Channel.QueueUnbind(queue, routingKey, exchange, nil)
	if err != nil {
		log.Fatal("Failed to unbind queue: %w", err)
	}
	return nil
}

func (rc *RabbitMQClient) Close()error {
	if rc.Channel != nil {
		if err := rc.Channel.Close(); err != nil {
			log.Fatal("Failed to close channel :%w", err)
		}
	}
	if rc.Connection != nil {
		if err := rc.Connection.Close(); err != nil {
			log.Fatal("Failed to close connection : %w", err)
		}
	}
	return nil
}
