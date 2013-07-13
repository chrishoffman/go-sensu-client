package sensu

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
)

type ClientConfig struct {
	Name           string   `json:"name"`
	Address        string   `json:"address"`
	Subscripations []string `json:"subscriptions"`
}

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
	Client   ClientConfig           `json:"client"`
	Rabbitmq RabbitmqConfig         `json:"rabbitmq"`
	data     *simplejson.Json
}

type ConfigData struct {
	data map[string]interface{}
}

func LoadConfigs(configFile string, configDirs []string) (*Config, error) {
	js, ferr := parseFile(configFile)
	if ferr != nil {
		log.Printf("Unable to open config file: %s", ferr)
	}

	for _, dir := range configDirs {
		files, derr := ioutil.ReadDir(dir)
		if derr != nil {
			log.Printf("Unable to open config directory: %s", derr)
		}

		for _, f := range files {
			jsd, err := parseFile(filepath.Join(dir, f.Name()))
			if err != nil {
				log.Printf("Could not load %s: %s", f.Name(), err)
			}

			err = js.Extend(jsd)
			if err != nil {
				log.Printf("Error merging configs: %s", err)
			}
		}
	}

	//Reencoding merged JSON to parse to concrete type
	mergedJson, err := json.Marshal(js.data)
	if err != nil {
		return nil, fmt.Errorf("Unable to reencode merged json")
	}
	config := new(Config)
	json.Unmarshal(mergedJson, &config)
	config.data, _ = simplejson.NewJson(mergedJson)

	return config, nil
}

func parseFile(filename string) (*ConfigData, error) {
	j := new(ConfigData)

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

func (j1 *ConfigData) Extend(j2 *ConfigData) (err error) {
	j1.data, err = mapExtend(j1.data, j2.data)
	return
}

func mapExtend(base map[string]interface{}, ext map[string]interface{}) (map[string]interface{}, error) {
	var err error
	for key, baseVal := range base {
		if extVal, ok := ext[key]; ok {
			b := reflect.ValueOf(baseVal)
			e := reflect.ValueOf(extVal)

			// Keep value of base if types do not match
			if b.Type() != e.Type() {
				return nil, fmt.Errorf("Conflicting types for key: %s (%s/%s). Skipping.", key, b.Kind().String(), e.Kind().String())
			}

			switch b.Kind() {
			case reflect.Slice:
				bSlice := (baseVal).([]interface{})
				eSlice := (extVal).([]interface{})

				for _, ele := range eSlice {
					bSlice = sliceExtend(bSlice, ele)
				}
				base[key] = bSlice
			case reflect.Map:
				bMap := (baseVal).(map[string]interface{})
				eMap := (extVal).(map[string]interface{})
				base[key], err = mapExtend(bMap, eMap)
				if err != nil {
					return nil, err
				}
			default:
				// Do nothing, prefer base
			}
		}
	}

	// Add all keys from ext that do not exist in base
	for key, extVal := range ext {
		if _, ok := base[key]; !ok {
			base[key] = extVal
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
