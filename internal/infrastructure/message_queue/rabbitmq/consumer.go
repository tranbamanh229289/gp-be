package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	client *RabbitMQClient
}

func NewConsumer(client *RabbitMQClient) *Consumer{
	return &Consumer{client: client}
}

func (c *Consumer) Consume(ctx context.Context, queue, tag string) error {
	msgs, err := c.client.Channel.Consume(queue, tag, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register consume: %v", err)
	}

	channelClosed := make(chan *amqp091.Error, 1)
	c.client.Channel.NotifyClose(channelClosed)

	go func(){
		defer log.Printf("Consumer goroutine exited")
		for {
			select {
				case msg, ok :=<-msgs: 
					if !ok {
						return
					}
					err := HandleMessage(ctx, msg)
					if err != nil {
						log.Printf("Failed to handle message: %v", err)
						if err := msg.Nack(false, true); err != nil {
							log.Printf("Failed to nack message: %v", err)
						}
						continue
					}
					
					if err := msg.Ack(false); err != nil {
						log.Printf("Failed to ack message: %v", err)
					} else {
						log.Printf("Message processed and acked: %s", msg.Body)
					}
				case <-ctx.Done():
					log.Println(" Context cancelled in goroutine !")
					c.Cancel(tag)
					return
				case err := <- channelClosed: 
					log.Printf("Channel closed: %v", err)
					return 
			}
		}

	}()

	select {
		case <- ctx.Done(): 
			log.Println("Consumer stopped by context")
			return ctx.Err()
		case err := <- channelClosed: 
			log.Printf("Channel closed: %w", err)
			return err

	}

}

func HandleMessage(ctx context.Context, msg amqp091.Delivery) error {
	log.Printf("Processing message %s", msg.Body)
	return nil
}

func (c *Consumer)Cancel(tag string) error {
	return c.client.Channel.Cancel(tag, false)
}