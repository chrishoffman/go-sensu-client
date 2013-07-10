package main

import (
	"flag"
	"log"
	"os"
	"sensu-client/sensu"
)

var configFile, configDir string

func init() {
	flag.StringVar(&configFile, "config-file", "config.json", "Default config file location")
	flag.StringVar(&configDir, "config-dir", "conf.d", "Default config directory to load config files")
	flag.Parse()
}

func main() {
	settings, err := sensu.LoadSettings(configFile, configDir)
	if err != nil {
		log.Printf("Unable to load settings: %s", err)
		os.Exit(1)
	}

	c := sensu.NewClient(settings)

	errc := make(chan error)
	c.Start(errc)
}
