package sensu

import (
	"github.com/streadway/amqp"
	"log"
)

type Subscriber struct {
	r MessageQueuer
	s []subscription
}

type subscription struct {
	s string
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
