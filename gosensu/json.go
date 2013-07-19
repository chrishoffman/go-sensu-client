package sensu

import (
	"fmt"
	"reflect"
)

type Json struct {
	data map[string]interface{}
}

func (j1 *Json) Extend(j2 *Json) (err error) {
	j1.data, err = mapExtend(j1.data, j2.data)
	return
}

func mapExtend(base map[string]interface{}, ext map[string]interface{}) (map[string]interface{}, error) {
	var err error
	for key, baseVal := range base {
		if extVal, ok := ext[key]; ok {
			b := reflect.ValueOf(baseVal)
			e := reflect.ValueOf(extVal)

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
