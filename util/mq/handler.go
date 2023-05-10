package mq

import (
	"errors"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MqHandler interface {
	Start() error
	Stop() error
	RegisterQueueListener(id string, listener func(*amqp.Delivery) error) error
	DeregisterQueueListener(id string) error
}

type mqHandler struct {
	mu         sync.RWMutex
	deliveryCh <-chan amqp.Delivery
	stopped    bool
	listeners  map[string]func(*amqp.Delivery) error
	doneCh     chan struct{}
	mapMu      sync.Mutex
}

func NewMqHandler(deliveryCh <-chan amqp.Delivery) MqHandler {
	return &mqHandler{
		deliveryCh: deliveryCh,
		doneCh:     make(chan struct{}, 1),
	}
}

func (h *mqHandler) Start() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.stopped {
		return nil
	}

	go func(doneCh chan struct{}, deliveryCh <-chan amqp.Delivery) error {
		for {
			select {
			case <-doneCh:
				return nil
			case payload, ok := <-deliveryCh:
				if ok {
					var wg sync.WaitGroup
					errCh := make(chan error)
					for _, listener := range h.listeners {
						wg.Add(1)
						// Worker goroutine
						go func(errCh chan<- error, listener func(*amqp.Delivery) error) {
							defer wg.Done()
							errCh <- listener(&payload)
						}(errCh, listener)
					}
					// Clean up goroutine
					go func(wg *sync.WaitGroup) {
						wg.Wait()
						close(errCh)
					}(&wg)
					// Check if any listeners failed to process delivery
					for err := range errCh {
						if err != nil {
							payload.Nack(false, false)
						}
					}
					payload.Ack(false)
				} else {
					return errors.New("delivery channel closed")
				}
			}
		}
	}(h.doneCh, h.deliveryCh)
	return nil
}

func (h *mqHandler) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.doneCh <- struct{}{}
	h.stopped = true
	return nil
}

func (h *mqHandler) RegisterQueueListener(consumer string, listener func(*amqp.Delivery) error) error {
	h.mapMu.Lock()
	defer h.mapMu.Unlock()
	h.Stop()
	defer h.Start()
	h.listeners[consumer] = listener
	return nil
}

func (h *mqHandler) DeregisterQueueListener(consumer string) error {
	h.mapMu.Lock()
	defer h.mapMu.Unlock()
	h.Stop()
	defer func() {
		if len(h.listeners) > 0 {
			h.Start()
		}
	}()
	delete(h.listeners, consumer)
	return nil
}
