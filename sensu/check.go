package sensu

type Check struct {
	Name            string   `json:"name"`
	Command         string   `json:""`
	Executed        int      `json:"executed"`
	Output          string   `json:"output"`
	Status          int      `json:"status"`
	Duration        float64  `json:"duration"`
	Timeout         int      `json:"timeout"`
	Handle          bool     `json:"handle"`
	Handlers        []string `json:"handlers"`
	commandExecuted string
}
