package sensu

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Subscriber struct {
	subs       []string
	deliveries <-chan amqp.Delivery
	done       chan error
}

func (s *Subscriber) Init(q MessageQueuer, c *Config) error {
	log.Printf("Declaring Queue")
	queue, err := q.QueueDeclare("")
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}
	log.Printf("declared Queue")

	for _, sub := range s.subs {
		log.Printf("declaring Exchange (%q)", sub)
		err = q.ExchangeDeclare(sub, "fanout")
		if err != nil {
			return fmt.Errorf("Exchange Declare: %s", err)
		}

		log.Printf("binding to Exchange %q", sub)
		err = q.QueueBind(queue.Name, "", sub)
		if err != nil {
			return fmt.Errorf("Queue Bind: %s", err)
		}
	}

	log.Printf("starting Consume")
	s.deliveries, err = q.Consume(queue.Name, "")
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	s.done = make(chan error)
	return nil
}

func (s *Subscriber) Start() {
	go handle(s.deliveries, s.done)

	// for {
	// 	select {
	// 	case <-s.done:
	// 		return
	// 	}
	// }
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
