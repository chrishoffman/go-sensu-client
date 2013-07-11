package sensu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

type RabbitmqConfigSSL struct {
	PrivateKeyFile string `json:"private_key_file"`
	CertChainFile  string `json:"cert_chain_file"`
}

type RabbitmqConfig struct {
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Vhost    string            `json:"vhost"`
	User     string            `json:"user"`
	Password string            `json:"password"`
	Ssl      RabbitmqConfigSSL `json:"ssl"`
}

type Config struct {
	Checks   map[string]interface{} `json:"checks"`
	Client   map[string]interface{} `json:"client"`
	Rabbitmq RabbitmqConfig         `json:"rabbitmq"`
	data     *Json
}

type Json struct {
	data map[string]interface{}
}

func NewJson(data map[string]interface{}) *Json {
	return &Json{data}
}

func LoadConfigs(configFile string, configDir string) (*Config, error) {
	js, ferr := parseFile(configFile)
	if ferr != nil {

	}

	files, derr := ioutil.ReadDir(configDir)
	if derr != nil {
		log.Printf("Unable to open config directory: %s", derr)
	}

	for _, f := range files {
		jsd, err := parseFile(f.Name())
		if err != nil {
			log.Printf("Could not load %s: %s", f, err)
		}

		err = js.Extend(jsd)
		if err != nil {
			log.Printf("Error merging configs: %s", err)
		}
	}

	//Reencoding merged JSON to parse to concrete type
	mergedJson, err := json.Marshal(js.data)
	if err != nil {
		return nil, fmt.Errorf("Unable to reencode merged json")
	}
	config := new(Config)
	json.Unmarshal(mergedJson, &config)
	config.data = js

	return config, nil
}

func parseFile(filename string) (*Json, error) {
	j := new(Json)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return j, fmt.Errorf("File error: %v", err)
	}

	err = json.Unmarshal(file, &j.data)
	if err != nil {
		return j, fmt.Errorf("json error: %v", err)
	}

	return j, nil
}

func (j1 *Json) Extend(j2 *Json) error {
	var err error
	j1.data, err = mapExtend(j1.data, j2.data)
	return err
}

func mapExtend(base map[string]interface{}, ext map[string]interface{}) (map[string]interface{}, error) {
	var err error
	for key, baseVal := range base {
		if extVal, ok := ext[key]; ok {
			b := reflect.ValueOf(baseVal)
			e := reflect.ValueOf(extVal)

			// Keep value of base if types do not match
			if b.Type() != e.Type() {
				return nil, fmt.Errorf("Conflicting types for key: %s (%s/%s)", key, b.Kind().String(), e.Kind().String())
			}

			switch b.Kind() {
			case reflect.Slice:
				bSlice := (baseVal).([]interface{})
				eSlice := (extVal).([]interface{})

				for _, ele := range eSlice {
					base[key] = sliceExtend(bSlice, ele)
				}
			case reflect.Map:
				bMap := (baseVal).(map[string]interface{})
				eMap := (extVal).(map[string]interface{})
				base[key], err = mapExtend(bMap, eMap)
				if err != nil {
					return nil, err
				}
			default:
				base[key] = extVal
			}
		}
	}
	return base, nil
}

func sliceExtend(slice []interface{}, i interface{}) []interface{} {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
