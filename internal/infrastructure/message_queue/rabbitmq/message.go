package rabbitmq

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Message struct {
	ID string
	Producer string
	Version string
	Priority uint8
	Timestamp time.Time

	RoutingKey string
	Headers amqp091.Table

	Payload any
}