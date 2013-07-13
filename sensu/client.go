package sensu

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Processor interface {
	Start(MessageQueuer, *Config)
	Stop()
}

type Client struct {
	config    *Config
	processes []Processor
}

func NewClient(c *Config, p []Processor) *Client {
	return &Client{
		config:    c,
		processes: p,
	}
}

func (c *Client) Start(errc chan error) {
	var disconnected chan *amqp.Error
	connected := make(chan bool)

	q := NewRabbitmq(c.config.Rabbitmq)
	go q.Connect(connected, errc)

	for {
		select {
		case <-connected:
			for _, proc := range c.processes {
				go proc.Start(q, c.config)
			}
			// Enable disconnect channel
			disconnected = q.Disconnected()
		case errd := <-disconnected:
			// Disable disconnect channel
			disconnected = nil

			log.Printf("RabbitMQ disconnected: %s", errd)
			for _, proc := range c.processes {
				proc.Stop()
			}

			time.Sleep(10 * time.Second)
			go q.Connect(connected, errc)
		}
	}
}
