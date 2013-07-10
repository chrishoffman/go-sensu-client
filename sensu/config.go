package sensu

import (
	// "encoding/json"
	"fmt"
	"reflect"
	// "log"
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
	data map[string]interface{}
}

// func Configure(configFile string, configDir string) (*Config, error) {

// }

// func parse(filename string) (*Config, error) {

// }

func extend(base map[string]interface{}, ext map[string]interface{}) (map[string]interface{}, error) {
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
					base[key] = appendUnique(bSlice, ele)
				}
			case reflect.Map:
				bMap := (baseVal).(map[string]interface{})
				eMap := (extVal).(map[string]interface{})
				if base[key], err = extend(bMap, eMap); err != nil {
					return nil, err
				}
			default:
				base[key] = extVal
			}
		}
	}
	return base, nil
}

func appendUnique(slice []interface{}, i interface{}) []interface{} {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}