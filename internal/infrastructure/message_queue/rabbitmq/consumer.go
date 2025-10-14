package rabbitmq

import (
	"be/pkg/logger"
	"context"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer struct {
	queue *RabbitQueue
	logger *logger.ZapLogger
}

func NewConsumer(queue *RabbitQueue, logger *logger.ZapLogger) *Consumer{
	return &Consumer{queue: queue, logger: logger}
}

func (c *Consumer) Consume(ctx context.Context, queue, tag string) error {
	msgs, err := c.queue.channel.Consume(queue, tag, false, false, false, false, nil)
	if err != nil {
		c.logger.Error("Failed to register consume: ", zap.Error(err))
		return err
	}

	channelClosed := make(chan *amqp091.Error, 1)
	c.queue.channel.NotifyClose(channelClosed)

	go func(){
		defer c.logger.Info("Consumer goroutine exited")
		for {
			select {
				case msg, ok :=<-msgs: 
					if !ok {
						return
					}
					err := HandleMessage(ctx, msg)
					if err != nil {
						c.logger.Info("Failed to handle message:", zap.Error(err))
						if err := msg.Nack(false, true); err != nil {
							c.logger.Info("Failed to nack message:", zap.Error(err))
						}
						continue
					}
					
					if err := msg.Ack(false); err != nil {
						c.logger.Info("Failed to ack message:", zap.Error(err))
					} else {
						c.logger.Info("Message processed and acked: ", zap.String("body", string(msg.Body)))
					}
				case <-ctx.Done():
					c.logger.Info(" Context cancelled in goroutine !")
					c.Cancel(tag)
					return
				case err := <- channelClosed: 
					c.logger.Error("Channel closed:", zap.Error(err))
					return
			}
		}

	}()

	select {
		case <- ctx.Done(): 
			c.logger.Info("Consumer stopped by context")
			return ctx.Err()
		case err := <- channelClosed: 
			c.logger.Error("Channel closed:", zap.Error(err))
			return err

	}

}

func HandleMessage(ctx context.Context, msg amqp091.Delivery) error {
	return nil
}

func (c *Consumer)Cancel(tag string) error {
	return c.queue.channel.Cancel(tag, false)
}