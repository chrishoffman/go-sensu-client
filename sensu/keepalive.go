package sensu

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)

type Keepalive struct {
	q        MessageQueuer
	config   *Config
	close    chan bool
}

const keepaliveInterval = 20 * time.Second

func NewKeepalive(q MessageQueuer, config *Config) *Keepalive {
	return &Keepalive{
		q:        q,
		config:   config,
		close:    make(chan bool),
	}
}

func (k *Keepalive) Start() {
	if err := k.q.ExchangeDeclare(
		"keepalives",
		"direct",
	); err != nil {
		log.Println("Exchange Declare: %s", err)
	}

	reset := make(chan bool)
	timer := time.AfterFunc(0, func() {
		k.publish(time.Now())
		reset <- true
	})
	defer timer.Stop()

	for {
		select {
		case <-reset:
			timer.Reset(keepaliveInterval)
		case <-k.close:
			return
		}
	}
}

func (k *Keepalive) Stop() {
	k.close <- true
}

func (k *Keepalive) publish(timestamp time.Time) {
	unixTimestamp := int64(timestamp.Unix())
	msg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         []byte(strconv.FormatInt(unixTimestamp, 10)),
		DeliveryMode: amqp.Persistent,
	}

	if err := k.q.Publish(
		"keepalives",
		"",
		msg,
	); err != nil {
		log.Printf("keepalive.publish: %v", err)
		return
	}
	log.Printf("Keepalive published: %s", strconv.FormatInt(unixTimestamp, 10))
}
