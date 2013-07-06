package main

import (
	"fmt"
	"flag"
	"github.com/bitly/go-simplejson"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/url"
	"log"
	"strconv"
	"time"
	"strings"
)

var configFile, configDir string

type SensuClient struct {
	ConfigFile string
	ConfigDir  string
	config     *simplejson.Json
	r          *rabbitmq
	k          *keepalive
}

type rabbitmq struct {
	uri          string
	conn         *amqp.Connection
	channel      *amqp.Channel
	disconnected chan *amqp.Error
}

type keepalive struct {
	interval time.Duration
	timer    *time.Timer
	restart  chan bool
	stop     chan bool
	close    chan bool
}

func (c *SensuClient) Start(errc chan error) {
	err := c.configure()
	if err != nil {
		errc <- fmt.Errorf("Unable to configure client")
	}

	connected := make(chan bool)
	e := make(chan error)
	c.r = new(rabbitmq)
	go c.r.Connect(c.config, connected, e)

	for {
		select {
		case <-connected:
			c.Keepalive(5 * time.Second)
		case <-c.r.disconnected:
			c.Reset()
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

	c.config = json
	return nil
}

func (c *SensuClient) Reset() chan error {

	return nil
}

func (c *SensuClient) Shutdown() chan error {

	return nil
}

func (r *rabbitmq) Connect(cfg *simplejson.Json, connected chan bool, errc chan error) {
	s, ok := cfg.CheckGet("rabbitmq"); 
	if !ok {
		errc <- fmt.Errorf("RabbitMQ settings missing from config")
		return
	}

	host := s.Get("host").MustString("localhost")
	port := int64(s.Get("port").MustInt(5672))
	user := s.Get("user").MustString("guest")
	password := s.Get("password").MustString("guest")
	vhost := s.Get("vhost").MustString("/sensu")

	userInfo := url.UserPassword(user, password)

	u := url.URL {
		Scheme: "amqp",
		Host: strings.Join([]string{host, ":", strconv.FormatInt(port, 10)}, ""),
		Path: vhost,
		User: userInfo,
	}
	uri := u.String()
	log.Println(uri)
	err := r.connect(uri)
	if err != nil {
		errc <- fmt.Errorf("Unable to connect to RabbitMQ")
		return
	}

	log.Println("RabbitMQ connected and channel established")
	connected <- true
}

func (r *rabbitmq) connect(uri string) error {
	var err error

	retryTicker := time.NewTicker(10 * time.Second)
	defer retryTicker.Stop()

	done := make(chan bool)
	go func() {
		for _ = range retryTicker.C {
			log.Printf("dialing %q", uri)
			r.conn, err = amqp.Dial(uri)
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

	return nil
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

	k.publish(r)
	reset := make(chan bool)
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

func NewClient(file string, dir string) *SensuClient {
	return &SensuClient{ConfigFile: configFile, ConfigDir: configDir}
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
	//time.Sleep(20 * time.Second)
	//panic("")

}
