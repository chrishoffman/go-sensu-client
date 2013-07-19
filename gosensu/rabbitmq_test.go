package sensu

import (
	"testing"
)

func Test_RabbitmqUrl(t *testing.T) {
	cfg := RabbitmqConfig{
		Host:     "localhost",
		Port:     5555,
		Vhost:    "/test",
		User:     "test",
		Password: "testpassword",
	}

	expected := "amqp://test:testpassword@localhost:5555//test"
	actual := createRabbitmqUri(cfg)
	if expected != actual {
		t.Errorf("[rmq url] expected: %s, actual: %s", expected, actual)
	}
}
