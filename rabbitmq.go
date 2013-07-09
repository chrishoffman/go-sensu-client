package main

import (
	"crypto/tls"
	"github.com/streadway/amqp"
	"log"
	"net/url"
	"strconv"
	"time"
)

type MessageQueuer interface {
	Connect(cfg RabbitmqConfig, connected chan bool, errc chan error)
	Disconnected() chan *amqp.Error
	ExchangeDeclare(name string, kind string) error	
}

type MessagePublisher interface {
	MessageQueuer
	Publish(exchange string, key string, msg amqp.Publishing) error
}

type MessageConsumer interface {
	MessageQueuer
}

type Rabbitmq struct {
	uri          string
	conn         *amqp.Connection
	channel      *amqp.Channel
	disconnected chan *amqp.Error
}

type RabbitmqConfig struct {
	Host     string
	Port     int
	Vhost    string
	User     string
	Password string
	Ssl      *tls.Config
}

const rabbitmqRetryInterval = 5 * time.Second

func (r *Rabbitmq) Connect(cfg RabbitmqConfig, connected chan bool, errc chan error) {

	u := url.URL{
		Scheme: "amqp",
		Host:   cfg.Host + ":" + strconv.FormatInt(int64(cfg.Port), 10),
		Path:   cfg.Vhost,
		User:   url.UserPassword(cfg.User, cfg.Password),
	}
	uri := u.String()

	reset := make(chan bool)
	done := make(chan bool)
	timer := time.AfterFunc(0, func() {
		r.connect(uri, done)
		reset <- true
	})
	defer timer.Stop()

	for {
		select {
		case <-done:
			log.Println("RabbitMQ connected and channel established")
			connected <- true
			return
		case <-reset:
			timer.Reset(rabbitmqRetryInterval)
		}
	}
}

func (r *Rabbitmq) Disconnected() chan *amqp.Error {
	return r.disconnected
}

func (r *Rabbitmq) ExchangeDeclare(name, kind string) error {
	return r.channel.ExchangeDeclare(
		name,
		kind,
		false, // All exchanges are not declared durable
		false,
		false,
		false,
		nil,
	)
}

func (r *Rabbitmq) Publish(exchange, key string, msg amqp.Publishing) error {
	return r.channel.Publish(
		exchange,
		key,
		false,
		false,
		msg,
	)
}

func (r *Rabbitmq) connect(uri string, done chan bool) {
	var err error

	log.Printf("dialing %q", uri)
	r.conn, err = amqp.Dial(uri)
	if err != nil {
		log.Printf("Dial: %s", err)
		return
	}

	log.Printf("Connection established, getting Channel")
	r.channel, err = r.conn.Channel()
	if err != nil {
		log.Printf("Channel: %s", err)
		return
	}

	// Notify disconnect channel when disconnected
	r.disconnected = make(chan *amqp.Error)
	r.channel.NotifyClose(r.disconnected)

	done <- true
}
