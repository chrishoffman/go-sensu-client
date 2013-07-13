package sensu

import (
	"github.com/bitly/go-simplejson"
	// "github.com/streadway/amqp"
	"testing"
	"time"
)

func Test_KeepalivePayload(t *testing.T) {
	timestamp := time.Now()
	config, _ := simplejson.NewJson([]byte(`{"name":"test","address":"1.2.3.4"}`))
	payload := createKeepalivePayload(config, timestamp)

	payloadBody, _ := simplejson.NewJson(payload.Body)
	if val, ok := payloadBody.CheckGet("timestamp"); ok {
		v, err := val.Int()
		if err != nil {
			t.Error("Unable to convert timestamp to int")
			return
		}

		bodyTs := time.Unix(int64(v), 0)
		if bodyTs != timestamp.Round(time.Second) {
			t.Errorf("timestamps do not match (%s/%s)", bodyTs, timestamp)
		}
	} else {
		t.Errorf("timestamp not found in payload body: %v", payloadBody)
	}

	if _, ok := payloadBody.CheckGet("name"); !ok {
		t.Error("Additional config not included in payload")
	}

}
