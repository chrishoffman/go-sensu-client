package sensu

import (
	"github.com/streadway/amqp"
	"github.com/bitly/go-simplejson"
	"log"
	"time"
)

type Keepalive struct {
	q      MessageQueuer
	config *Config
	close  chan bool
}

const keepaliveInterval = 20 * time.Second

func (k *Keepalive) Init(q MessageQueuer, config *Config) {
	k.q = q
	k.config = config
	k.close = make(chan bool)
}

func (k *Keepalive) Start() {
	if err := k.q.ExchangeDeclare(
		"keepalives",
		"direct",
	); err != nil {
		log.Println("Exchange Declare: %s", err)
		//panic?
	}

	clientConfig := k.config.Data().Get("client")
	reset := make(chan bool)
	timer := time.AfterFunc(0, func() {
		payload := createKeepalivePayload(clientConfig, time.Now())
		k.publish(payload)
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

func (k *Keepalive) publish(payload amqp.Publishing) {
	if err := k.q.Publish(
		"keepalives",
		"",
		payload,
	); err != nil {
		log.Printf("keepalive.publish: %v", err)
		return
	}
	log.Print("Keepalive published")
}

func createKeepalivePayload(clientConfig *simplejson.Json, timestamp time.Time) (amqp.Publishing) {
	payload := clientConfig
	payload.Set("timestamp", int64(timestamp.Unix()))
	body, _ := payload.MarshalJSON()
	return amqp.Publishing{
		ContentType:  "application/octet-stream",
		Body:         body,
		DeliveryMode: amqp.Transient,
	}
}
