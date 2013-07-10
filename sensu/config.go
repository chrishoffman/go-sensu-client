package sensu

import (
	// "encoding/json"
	// "log"
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
	data map[string]interface{}
}

// func Configure(configFile string, configDir string) (*Config, error) {

// }

func (c1 *Config) Extend(c2 *Config) error {
	var err error
	c1.data, err = mapExtend(c1.data, c2.data)
	return err
}

// func parse(filename string) (*Config, error) {

// }
