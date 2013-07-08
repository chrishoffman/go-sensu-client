package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/streadway/amqp"
	"log"
	"net/url"
	"strconv"
	"time"
)

type MessageQueuer interface {
	// TODO: Fixing signature for config, needs to be just rmq settings
	Connect(*simplejson.Json, chan bool, chan error)
	ExchangeDeclare(string, string, bool, bool, bool, bool, amqp.Table) error
}

type Rabbitmq struct {
	uri          string
	conn         *amqp.Connection
	channel      *amqp.Channel
	disconnected chan *amqp.Error
}

type RabbitmqConfig struct {
	host     string
	port     int64
	vhost    string
	user     string
	password string
	ssl      *tls.Config
}

const rabbitmqRetryInterval = 5 * time.Second

func (r *Rabbitmq) Connect(cfg *simplejson.Json, connected chan bool, errc chan error) {
	s, ok := cfg.CheckGet("rabbitmq")
	if !ok {
		errc <- fmt.Errorf("RabbitMQ settings missing from config")
		return
	}

	host := s.Get("host").MustString()
	port := s.Get("port").MustInt()
	user := s.Get("user").MustString()
	password := s.Get("password").MustString()
	vhost := s.Get("vhost").MustString()

	userInfo := url.UserPassword(user, password)

	u := url.URL{
		Scheme: "amqp",
		Host:   host + ":" + strconv.FormatInt(int64(port), 10),
		Path:   vhost,
		User:   userInfo,
	}
	uri := u.String()

	reset := make(chan bool)
	done := make(chan bool)
	go r.connect(uri, done)
	timer := time.AfterFunc(rabbitmqRetryInterval, func() {
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
