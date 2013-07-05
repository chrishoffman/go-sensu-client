package main

import (
	"flag"
	"github.com/bitly/go-simplejson"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

var configFile, configDir string


type SensuClient struct {
	ConfigFile string
	ConfigDir  string
	settings   *simplejson.Json
	r          *rabbitmq
	k          *keepalive
}

type rabbitmq struct {
	uri          string
	conn         *amqp.Connection
	channel      *amqp.Channel
	connected    chan bool
	disconnected chan *amqp.Error
}

type keepalive struct {
	interval time.Duration
	timer    *time.Timer
	restart  chan bool
	stop     chan bool
	close    chan bool
}

func (c *SensuClient) Start(err chan error) {
	c.configure()

	c.r = &rabbitmq{
		uri:          "amqp://guest:guest@localhost:5672/",
		connected:    make(chan bool),
		disconnected: make(chan *amqp.Error),
	}

	go c.r.Connect()

	for {
		select {
		case <-c.r.connected:
			c.Keepalive(5 * time.Second)
			// TODO: Implement disconnect channel
		}
	}
}

func (c *SensuClient) configure() error {
	file, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
	    log.Println("File error: %v", err)
	}

	json, err := simplejson.NewJson(file)
	if err != nil {
	    log.Println("json error: %v\n", err)
	}

	c.settings = json
	return nil
}

func (c *SensuClient) Reset() chan error {

	return nil
}

func (c *SensuClient) Shutdown() chan error {

	return nil
}

func (r *rabbitmq) Connect() {
	var err error

	retryTicker := time.NewTicker(10 * time.Second)
	defer retryTicker.Stop()

	done := make(chan bool)
	go func() {
		for _ = range retryTicker.C {
			log.Printf("dialing %q", r.uri)
			r.conn, err = amqp.Dial(r.uri)
			if err != nil {
				log.Println("Dial: %s", err)
				continue
			}

			log.Printf("got Connection, getting Channel")
			r.channel, err = r.conn.Channel()
			if err != nil {
				log.Println("Channel: %s", err)
				continue
			}

			done <- true
		}
	}()
	<-done

	// Notify disconnect channel when disconnected
	r.channel.NotifyClose(r.disconnected)

	log.Println("RabbitMQ connected and channel established")
	r.connected <- true
}

func (c *SensuClient) Keepalive(interval time.Duration) {
	c.k = &keepalive{
		interval: interval,
		restart:  make(chan bool),
		stop:     make(chan bool),
		close:    make(chan bool),
	}

	go c.k.loop(c.r)
}

func (k *keepalive) loop(r *rabbitmq) {
	if err := r.channel.ExchangeDeclare(
		"keepalives",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Println("Exchange Declare: %s", err)
	}

	reset := make(chan bool)
	k.publish(r)
	k.timer = time.AfterFunc(k.interval, func() {
		k.publish(r)
		reset <- true
	})

	for {
		select {
		case <-reset:
			k.timer.Reset(k.interval)
		case <-k.restart:
			k.timer.Reset(k.interval)
		case <-k.stop:
			k.timer.Stop()
		case <-k.close:
			k.timer.Stop()
			return
		}
	}
}

func (k *keepalive) publish(r *rabbitmq) {
	unixTimestamp := int64(time.Now().Unix())
	msg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         []byte(strconv.FormatInt(unixTimestamp, 10)),
		DeliveryMode: amqp.Persistent,
	}

	if err := r.channel.Publish(
		"keepalives",
		"",
		false,
		false,
		msg,
	); err != nil {
		log.Printf("keepalive.publish: %v", err)
		return
	}
	log.Printf("Keepalive published: %s", strconv.FormatInt(unixTimestamp, 10))
}

func init() {
	flag.StringVar(&configFile, "config-file", "config.json", "Default config file location")
	flag.StringVar(&configDir, "config-dir", "conf.d", "Default config directory to load config files")
	flag.Parse()
}

func main() {

	c := &SensuClient{
		ConfigFile: configFile,
		ConfigDir:  configDir,
	}

	e := make(chan error)
	go c.Start(e)

	for {
		select {
		case err := <-e:
			panic(err)
		}
	}
	//time.Sleep(20 * time.Second)
	//panic("")

}
