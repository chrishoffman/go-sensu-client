package sensu

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
