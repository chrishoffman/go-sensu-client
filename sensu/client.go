package sensu

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Client struct {
	config *Config
	r      MessageQueuer
	k      *Keepalive
}

func NewClient(s *Config) *Client {
	return &Client{
		r:      new(Rabbitmq),
		config: s,
	}
}

func (c *Client) Start(errc chan error) {
	var disconnected chan *amqp.Error

	connected := make(chan bool)
	go c.r.Connect(c.config.Rabbitmq, connected, errc)

	for {
		select {
		case <-connected:
			c.Keepalive(5 * time.Second)
			disconnected = c.r.Disconnected() // Enable disconnect channel
		case errd := <-disconnected:
			log.Printf("RabbitMQ disconnected: %s", errd)
			c.Reset()
			disconnected = nil // Disable disconnect channel
			time.Sleep(10 * time.Second)
			go c.r.Connect(c.config.Rabbitmq, connected, errc)
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
	c.k = NewKeepalive(c.r, interval)
	go c.k.Start()
}
