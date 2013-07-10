package sensu

import (
	// "encoding/json"
	"fmt"
	"reflect"
)

type RabbitmqConfigSSL struct {
	PrivateKeyFile string
	CertChainFile  string
}

type RabbitmqConfig struct {
	Host     string
	Port     int
	Vhost    string
	User     string
	Password string
	Ssl      RabbitmqConfigSSL
}

type Config struct {
	d map[string]interface{}
}

// func Configure(configFile string, configDir string) (*Config, error) {

// }

// func parse(filename string) (*Config, error) {

// }

func (base *Config) merge(ext *Config) error {
	for key, _ := range base.d {
		b := reflect.ValueOf(base.d[key])
		e := reflect.ValueOf(ext.d[key])

		// Keep value of base if types do not match
		if b.Type() != e.Type() {
			return fmt.Errorf("Conflicting types for key: %s (%s/%s)", key, b.Kind().String(), e.Kind().String())
		}

		switch b.Kind() {
		case reflect.Slice:
			continue
		case reflect.Map:
			continue
		default:
			continue
		} 
	}
	return nil
}
