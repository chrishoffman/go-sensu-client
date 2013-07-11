package sensu

import (
// "github.com/streadway/amqp"
// "log"
// "strconv"
// "time"
)

type Subscriber struct {
	r MessageQueuer
	s []subscription
}

type subscription struct {
	s string
}
