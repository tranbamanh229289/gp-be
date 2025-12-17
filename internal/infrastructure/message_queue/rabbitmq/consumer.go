package rabbitmq

import (
	"be/pkg/logger"
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	queue  *RabbitQueue
	logger *logger.ZapLogger
}

func NewConsumer(queue *RabbitQueue, logger *logger.ZapLogger) (*Consumer, error) {
	return &Consumer{queue: queue, logger: logger}, nil
}

func (c *Consumer) Consume(ctx context.Context, queue, tag string) error {
	msgs, err := c.queue.channel.Consume(queue, tag, false, false, false, false, nil)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Failed to register consume: %s", err))
		return err
	}

	channelClosed := make(chan *amqp091.Error, 1)
	c.queue.channel.NotifyClose(channelClosed)

	go func() {
		defer c.logger.Info("Consumer goroutine exited")
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				err := HandleMessage(ctx, msg)
				if err != nil {
					c.logger.Error(fmt.Sprintf("Failed to handle message: %s", err))
					if err := msg.Nack(false, true); err != nil {
						c.logger.Error(fmt.Sprintf("Failed to nack message: %s", err))
					}
					continue
				}

				if err := msg.Ack(false); err != nil {
					c.logger.Info(fmt.Sprintf("Failed to ack message: %s", err))
				} else {
					c.logger.Info(fmt.Sprintf("Message processed and acked: %s", msg.Body))
				}
			case <-ctx.Done():
				c.logger.Info("Context cancelled in goroutine !")
				c.Cancel(tag)
				return
			case err := <-channelClosed:
				c.logger.Error(fmt.Sprintf("Channel closed: %s", err))
				return
			}
		}

	}()

	select {
	case <-ctx.Done():
		c.logger.Info("Consumer stopped by context")
		return ctx.Err()
	case err := <-channelClosed:
		c.logger.Error(fmt.Sprintf("Channel closed: %s", err))
		return err

	}

}

func HandleMessage(ctx context.Context, msg amqp091.Delivery) error {
	return nil
}

func (c *Consumer) Cancel(tag string) error {
	return c.queue.channel.Cancel(tag, false)
}
