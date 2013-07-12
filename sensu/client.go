package sensu

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Processor interface {
	Start()
	Stop()
}

type Client struct {
	config    *Config
	q         MessageQueuer
	processes []Processor
}

func NewClient(c *Config) *Client {
	return &Client{
		config:  c,
	}
}

func (c *Client) Start(errc chan error) {
	var disconnected chan *amqp.Error
	connected := make(chan bool)

	c.q = NewRabbitmq(c.config.Rabbitmq)
	go c.q.Connect(connected, errc)

	c.processes = []Processor{NewKeepalive(c.q, 5 * time.Second)}


	for {
		select {
		case <-connected:
			for _, proc := range c.processes {
				go proc.Start()
			}
			// Enable disconnect channel
			disconnected = c.q.Disconnected()
		case errd := <-disconnected:
			// Disable disconnect channel
			disconnected = nil

			log.Printf("RabbitMQ disconnected: %s", errd)
			for _, proc := range c.processes {
				proc.Stop()
			}

			time.Sleep(10 * time.Second)
			go c.q.Connect(connected, errc)
		}
	}
}
