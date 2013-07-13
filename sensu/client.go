package sensu

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Processor interface {
	Init(MessageQueuer, *Config)
	Start()
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

func (c *Client) Start() {
	var disconnected chan *amqp.Error
	connected := make(chan bool)

	q := NewRabbitmq(c.config.Rabbitmq)
	go q.Connect(connected)

	for {
		select {
		case <-connected:
			for _, proc := range c.processes {
				proc.Init(q, c.config)
				go proc.Start()
			}
			// Enable disconnect channel
			disconnected = q.Disconnected()
		case errd := <-disconnected:
			// Disable disconnect channel
			disconnected = nil

			log.Printf("RabbitMQ disconnected: %s", errd)
			c.Stop()

			time.Sleep(10 * time.Second)
			go q.Connect(connected)
		}
	}
}

func (c *Client) Stop() {
	for _, proc := range c.processes {
		proc.Stop()
	}
}

func (c *Client) Shutdown() {
	// Disconnect rabbitmq
}
