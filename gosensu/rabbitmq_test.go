package sensu

// import (
// 	"strconv"
// 	"testing"
// )

// func Test_RabbitmqUrl(t *testing.T) {
// 	cfg := RabbitmqConfig{
// 		Host: "localhost",
// 		Port: 
// 	}


// }


// func createRabbitmqUri(cfg RabbitmqConfig) string {
// 	u := url.URL{
// 		Scheme: "amqp",
// 		Host:   fmt.Sprintf("%s:%s", cfg.Host, strconv.FormatInt(int64(cfg.Port), 10)),
// 		Path:   fmt.Sprintf("/%s", cfg.Vhost),
// 		User:   url.UserPassword(cfg.User, cfg.Password),
// 	}
// 	return u.String()
// }
