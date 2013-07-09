package sensu

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"time"
)

type Client struct {
	configFile string
	configDir  string
	config     *simplejson.Json
	r          *Rabbitmq
	k          *Keepalive
}


func NewClient(file string, dir string) *Client {
	return &Client{
		configFile: file,
		configDir:  dir,
		r:          new(Rabbitmq),
	}
}

func (c *Client) Start(errc chan error) {
	var disconnected chan *amqp.Error

	err := c.configure()
	if err != nil {
		errc <- fmt.Errorf("Unable to configure client")
		return
	}

	// Get RabbitMQ configs
	s, ok := c.config.CheckGet("rabbitmq")
	if !ok {
		errc <- fmt.Errorf("RabbitMQ settings missing from config")
		return
	}

	rmqConfig := RabbitmqConfig{
		Host:     s.Get("host").MustString(),
		Port:     s.Get("port").MustInt(),
		Vhost:    s.Get("vhost").MustString(),
		User:     s.Get("user").MustString(),
		Password: s.Get("password").MustString(),
	}

	connected := make(chan bool)
	go c.r.Connect(rmqConfig, connected, errc)

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
			go c.r.Connect(rmqConfig, connected, errc)
		}
	}
}

func (c *Client) configure() error {
	file, err := ioutil.ReadFile(c.configFile)
	if err != nil {
		log.Printf("File error: %v", err)
	}

	json, err := simplejson.NewJson(file)
	if err != nil {
		log.Printf("json error: %v\n", err)
	}

	c.config = json
	return nil
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
