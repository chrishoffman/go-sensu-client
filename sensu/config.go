package sensu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

type RabbitmqConfigSSL struct {
	PrivateKeyFile string
	CertChainFile  string
}

type RabbitmqConfig struct {
	Host     string
	Port     int
	Vhost    string
	User     string
	Password string
	Ssl      RabbitmqConfigSSL
}

type Config struct {
	Checks   map[string]interface{}
	Client   map[string]interface{}
	Rabbitmq RabbitmqConfig
}

type Json struct {
	data map[string]interface{}
}

func LoadConfigs(configFile string, configDir string) (*Config, error) {
	js, ferr := parseFile(configFile)
	if ferr != nil {

	}

	files, derr := ioutil.ReadDir(configDir)
	if derr != nil {
		//TODO: errors
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

	config := new(Config)

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
