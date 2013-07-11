package sensu_test

import (
	"fmt"
	. "sensu-client/sensu"
	"testing"
)

func Test_Extend(t *testing.T) {
	d1 := make(map[string]interface{})
	d2 := make(map[string]interface{})

	d1["bbb"] = "blah"
	d1["abc"] = []interface{}{"a", "b", "c"}
	d2["abc"] = []interface{}{"a", "d"}

	j1 := NewJson(d1)
	j2 := NewJson(d2)

	err := j1.Extend(j2)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", j1))
	}
}

// func Test_Parse(t *testing.T) {
// 	data, _ := parseFile("../config.json")
// 	t.Log(data)
// }
