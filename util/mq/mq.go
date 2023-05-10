package mq

import (
	"context"
	"fmt"
	"sync"

	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
	"github.com/rabbitmq/amqp091-go"
)

type MQ interface {
	RegisterQueueListener(queue string, id string, listener func(*amqp091.Delivery) error) error
	DeregisterQueueListener(queue string, id string) error
	SendToExchange(ctx context.Context, exchangeName string, routingKey string, payload Payload) error
}

type Payload interface {
	Serialize() (string, error)
	ContentType() string
	ContentEncoding() string
}

type mq struct {
	channel  *amqp091.Channel
	handlers map[string]MqHandler
	mu       sync.RWMutex
}

func NewRabbitMQ(consumer string, url string, username string, password string, port int) (MQ, error) {
	if username == "" {
		username = "guest"
	}
	if password == "" {
		password = "guest"
	}
	if url == "" {
		url = "localhost"
	}
	if port == 0 {
		port = 5672
	}
	return h.FlatMap2(
		h.FactoryM(func() (*amqp091.Connection, error) {
			amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, url, port)
			config := amqp091.Config{Properties: amqp091.NewConnectionProperties()}
			config.Properties.SetClientConnectionName(consumer)
			return amqp091.DialConfig(amqpURI, config)
		}),
		h.Lift(func(con *amqp091.Connection) (*amqp091.Channel, error) {
			return con.Channel()
		}),
		h.LiftJ(func(ch *amqp091.Channel) MQ {
			return &mq{channel: ch}
		}),
	).Eval()
}

func (q *mq) RegisterQueueListener(queue string, id string, listener func(*amqp091.Delivery) error) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if handler, ok := q.handlers[queue]; ok {
		return handler.RegisterQueueListener(id, listener)
	} else {
		if deliveryCh, err := q.channel.Consume(queue, id, true, false, false, false, nil); err != nil {
			return err
		} else {
			handler = NewMqHandler(deliveryCh)
			q.handlers[queue] = handler
			return handler.RegisterQueueListener(id, listener)
		}
	}
}

func (q *mq) DeregisterQueueListener(queue string, id string) error {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if handler, ok := q.handlers[queue]; ok {
		return handler.DeregisterQueueListener(id)
	}
	return nil
}

func (q *mq) SendToExchange(ctx context.Context, exchangeName string, routingKey string, payload Payload) error {
	return h.FlatMap(
		h.FactoryM(payload.Serialize),
		h.LiftE(func(data string) error {
			return q.channel.PublishWithContext(ctx, exchangeName, routingKey, true, false, amqp091.Publishing{
				Headers:         amqp091.Table{},
				ContentType:     payload.ContentType(),
				ContentEncoding: payload.ContentEncoding(),
				DeliveryMode:    amqp091.Persistent,
				Priority:        0,
				Body:            []byte(data),
			})
		}),
	).Error()
}
