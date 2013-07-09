package main

import (
	"flag"
	"sensu-client/sensu"
)

var configFile, configDir string

func init() {
	flag.StringVar(&configFile, "config-file", "config.json", "Default config file location")
	flag.StringVar(&configDir, "config-dir", "conf.d", "Default config directory to load config files")
	flag.Parse()
}

func main() {
	c := sensu.NewClient(configFile, configDir)

	errc := make(chan error)
	c.Start(errc)
}
