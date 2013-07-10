package sensu

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
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

func parse(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("File error: %v", err)
	}

	c := new(Config)
	err = json.Unmarshal(file, &c.data)
	if err != nil {
		return nil, fmt.Errorf("json error: %v\n", err)
	}

	return c, nil
}
