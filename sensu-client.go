package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"time"
)

var configFile, configDir string

type SensuClient struct {
	ConfigFile string
	ConfigDir  string
	config     *simplejson.Json
	r          *rabbitmq
	k          *Keepalive
}

func (c *SensuClient) Start(errc chan error) {
	var disconnected chan *amqp.Error

	err := c.configure()
	if err != nil {
		errc <- fmt.Errorf("Unable to configure client")
		return
	}

	connected := make(chan bool)
	go c.r.Connect(c.config, connected, errc)

	for {
		select {
		case <-connected:
			c.Keepalive(5 * time.Second)
			disconnected = c.r.disconnected // Enable disconnect channel
		case errd := <-disconnected:
			log.Printf("RabbitMQ disconnected: %s", errd)
			c.Reset()
			disconnected = nil // Disable disconnect channel
			time.Sleep(rabbitmqRetryInterval)
			go c.r.Connect(c.config, connected, errc)
		}
	}
}

func (c *SensuClient) configure() error {
	file, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		log.Printf("File error: %v", err)
	}

	json, err := simplejson.NewJson(file)
	if err != nil {
		log.Printf("json error: %v\n", err)
	}

	c.config = json
	return nil
}

func (c *SensuClient) Reset() chan error {
	// Stop keepalive timer
	c.k.Stop()

	return nil
}

func (c *SensuClient) Shutdown() chan error {

	return nil
}

func (c *SensuClient) Keepalive(interval time.Duration) {
	c.k = NewKeepalive(c.r, interval)
	go c.k.Start()
}

func NewClient(file string, dir string) *SensuClient {
	return &SensuClient{
		ConfigFile: configFile,
		ConfigDir:  configDir,
		r:          new(rabbitmq),
	}
}

func init() {
	flag.StringVar(&configFile, "config-file", "config.json", "Default config file location")
	flag.StringVar(&configDir, "config-dir", "conf.d", "Default config directory to load config files")
	flag.Parse()
}

func main() {

	c := NewClient(configFile, configDir)

	e := make(chan error)
	go c.Start(e)

	for {
		select {
		case err := <-e:
			panic(err)
		}
	}
}
