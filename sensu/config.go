package sensu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func LoadSettings(configFile string, configDir string) (*Config, error) {
	config, ferr := loadFile(configFile)
	if ferr != nil {

	}

	files, derr := ioutil.ReadDir(configDir)
	if derr != nil {
		//TODO: errors
	}

	for _, f := range files {
		cfg, err := loadFile(f.Name())
		if err != nil {
			log.Printf("Could not load %s: %s", f, ferr)
		}

		err = config.Extend(cfg)
		if err != nil {
			log.Printf("Error merging configs: %s", err)
		}
	}

	return config, nil
}

func (c *Config) Get(key string) (interface{}, bool) {
	if val, ok := c.data[key]; ok {
		return val, true
	}
	return nil, false
}

func (c *Config) Validate() (bool, error) {

	return true, nil
}

func (c1 *Config) Extend(c2 *Config) error {
	var err error
	c1.data, err = mapExtend(c1.data, c2.data)
	return err
}

func loadFile(filename string) (*Config, error) {
	c := new(Config)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, fmt.Errorf("File error: %v", err)
	}

	err = json.Unmarshal(file, &c.data)
	if err != nil {
		return c, fmt.Errorf("json error: %v", err)
	}

	return c, nil
}
