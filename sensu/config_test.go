package sensu_test

import (
	"fmt"
	"testing"
)

func Test_Extend(t *testing.T) {
	d1 := make(map[string]interface{})
	d2 := make(map[string]interface{})

	d1["bbb"] = "blah"
	d1["abc"] = []interface{}{"a", "b", "c"}
	d2["abc"] = []interface{}{"a", "d"}

	c1 := NewConfig(d1)
	c2 := NewConfig(d2)

	c1.Extend(c2)
	t.Log(fmt.Sprintf("%+v", c1))
}

func Test_Parse(t *testing.T) {
	data, _ := parse("../config.json")
	t.Log(data)
}
