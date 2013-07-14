package sensu

import (
	"github.com/streadway/amqp"
	"log"
)

type Subscriber struct {
	q    MessageQueuer
	s    []string
	done chan error
}

func (s *Subscriber) Init(q MessageQueuer, c *Config) {
	s.q = q
	s.done = make(chan error)
}

func (s *Subscriber) Start() {
	log.Printf("Declaring Queue")
	queue, err := s.q.QueueDeclare("")
	if err != nil {
		log.Printf("Queue Declare: %s", err)
	}
	log.Printf("declared Queue")

	for _, sub := range s.s {
		log.Printf("declaring Exchange (%q)", sub)
		err = s.q.ExchangeDeclare(sub, "fanout")
		if err != nil {
			log.Printf("Exchange Declare: %s", err)
		}

		log.Printf("binding to Exchange %q", sub)
		err = s.q.QueueBind(queue.Name, "", sub)
		if err != nil {
			log.Printf("Queue Bind: %s", err)
		}
	}

	log.Printf("starting Consume")
	deliveries, err := s.q.Consume(queue.Name, "")
	if err != nil {
		log.Printf("Queue Consume: %s", err)
	}

	go handle(deliveries, s.done)
	<-s.done
}

func (s *Subscriber) Stop() {
	s.done <- nil
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
