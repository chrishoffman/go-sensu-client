package sensu

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Processor interface {
	Start()
	Stop()
	Restart()
	Close()
}

type Client struct {
	config    *Config
	q         MessageQueuer
	processes []*Processor
	k         *Keepalive
}

func NewClient(c *Config) *Client {
	return &Client{
		config:  c,
	}
}

func (c *Client) Start(errc chan error) {
	var disconnected chan *amqp.Error

	c.q = NewRabbitmq(c.config.Rabbitmq)

	connected := make(chan bool)
	go c.q.Connect(connected, errc)

	for {
		select {
		case <-connected:
			c.Keepalive(5 * time.Second)
			disconnected = c.q.Disconnected() // Enable disconnect channel
		case errd := <-disconnected:
			log.Printf("RabbitMQ disconnected: %s", errd)
			c.Reset()
			disconnected = nil // Disable disconnect channel
			time.Sleep(10 * time.Second)
			go c.q.Connect(connected, errc)
		}
	}
}

func (c *Client) Reset() chan error {
	// Stop keepalive timer
	c.k.Stop()

	return nil
}

func (c *Client) Shutdown() chan error {

	return nil
}

func (c *Client) Keepalive(interval time.Duration) {
	c.k = NewKeepalive(c.q, interval)
	go c.k.Start()
}
