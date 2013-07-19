package sensu

import (
	"github.com/bitly/go-simplejson"
)

type Check struct {
	Name            string 
	Command         string
	Executed        int
	Status          int
	Output          string
	Duration        float64
	Timeout         int
	commandExecuted string
	data            *simplejson.Json
}
