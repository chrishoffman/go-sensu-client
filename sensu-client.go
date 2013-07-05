package main

import (
	_ "fmt"
	"log"
	"time"
	"strconv"
	"github.com/streadway/amqp"
)

type SensuClient struct {
	r *rabbitmq
	k *keepalive
}

type rabbitmq struct {
	uri     string
	conn    *amqp.Connection
	channel *amqp.Channel
	connected chan bool
	disconnected chan *amqp.Error
}

type keepalive struct {
	interval time.Duration
	timer    *time.Timer
	restart  chan bool
	stop     chan bool
	close    chan bool
}

func (c *SensuClient) Start() chan error {

    go c.r.Connect()

    for {
        select {
            case <- c.r.connected:
            	c.keepalive(5 * time.Second)
		}
	}
}

func (r *rabbitmq) Connect() {
    var err error

    retryTicker := time.NewTicker(10 * time.Second)
    defer retryTicker.Stop()

    done := make(chan bool)
    go func() {
        for _ = range retryTicker.C {
            log.Printf("dialing %q", r.uri)
            r.conn, err = amqp.Dial(r.uri)
            if err != nil {
                log.Println("Dial: %s", err)
                continue
            }

            log.Printf("got Connection, getting Channel")
            r.channel, err = r.conn.Channel()
            if err != nil {
                log.Println("Channel: %s", err)
                continue
            }

            done <- true
        }
    }()
    <-done

    // Notify disconnect channel when disconnected
    r.channel.NotifyClose(r.disconnected)

    log.Println("RabbitMQ connected and channel established")
    r.connected <- true
}

func (c *SensuClient) keepalive(interval time.Duration) {
	c.k = &keepalive{
		interval: interval,
		restart:  make(chan bool),
		stop:     make(chan bool),
		close:    make(chan bool),
	}

	go c.k.loop(c.r)
}

func (k *keepalive) loop(r *rabbitmq) {
    if err := r.channel.ExchangeDeclare(
        "keepalives",     
        "direct",
        true,
        false,
        false,
        false,
        nil,
    ); err != nil {
        log.Println("Exchange Declare: %s", err)
    }

	reset := make(chan bool)
	k.publish(r)
	k.timer = time.AfterFunc(k.interval, func() {
		k.publish(r)
		reset <- true
	})

	for {
		select {
		case <-reset:
			k.timer.Reset(k.interval)
		case <-k.restart:
			k.timer.Reset(k.interval)
		case <-k.stop:
			k.timer.Stop()
		case <-k.close:
			k.timer.Stop()
			return
		}
	}
}

func (k *keepalive) publish(r *rabbitmq) {
    unixTimestamp := int64(time.Now().Unix())
    msg := amqp.Publishing{
        ContentType:  "application/json",
        Body:         []byte(strconv.FormatInt(unixTimestamp, 10)),
        DeliveryMode: amqp.Persistent,
    }

    if err := r.channel.Publish(
        "keepalives",
        "",
        false,
        false,
        msg,
    ); err != nil {
        log.Printf("keepalive.publish: %v", err)
        return
    }
    log.Printf("Keepalive published: %s", strconv.FormatInt(unixTimestamp, 10))
}

func main() {
	r := &rabbitmq {
		uri: "amqp://guest:guest@localhost:5672/",
		connected:  make(chan bool),
		disconnected:     make(chan *amqp.Error),		
	}

	c := &SensuClient {
		r: r,
	}

	go c.Start()

	select {}
	//time.Sleep(20 * time.Second)
	//panic("")

}
