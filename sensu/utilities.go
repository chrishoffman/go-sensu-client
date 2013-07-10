package sensu

import (
	"fmt"
	"reflect"
)

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
					base[key] = sliceAppendUnique(bSlice, ele)
				}
			case reflect.Map:
				bMap := (baseVal).(map[string]interface{})
				eMap := (extVal).(map[string]interface{})
				if base[key], err = mapExtend(bMap, eMap); err != nil {
					return nil, err
				}
			default:
				base[key] = extVal
			}
		}
	}
	return base, nil
}

func sliceAppendUnique(slice []interface{}, i interface{}) []interface{} {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
