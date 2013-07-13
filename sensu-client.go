package main

import (
	"flag"
	"log"
	"os"
	"sensu-client/sensu"
	"strings"
)

var configFile, configDir string

func init() {
	flag.StringVar(&configFile, "config-file", "config.json", "Sensu JSON config file")
	flag.StringVar(&configDir, "config-dir", "conf.d", "directory or comma-delimited directory list for Sensu JSON config files")
	flag.Parse()
}

func main() {
	configDirs := strings.Split(configDir, ",")
	settings, err := sensu.LoadConfigs(configFile, configDirs)
	if err != nil {
		log.Printf("Unable to load settings: %s", err)
		os.Exit(1)
	}

	c := sensu.NewClient(settings)

	errc := make(chan error)
	c.Start(errc)
}
