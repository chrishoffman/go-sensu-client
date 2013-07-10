package sensu

import (
	"testing"
)

func Test_Extend(t *testing.T) {
	d1 := make(map[string]interface{})
	d2 := make(map[string]interface{})

	d1["bbb"] = "blah"
	d1["abc"] = []interface{}{"a","b","c"}
	d2["abc"] = []interface{}{"a","d"}

	c1 := &Config{d1}
	c2 := &Config{d2}

	c1.Extend(c2)
	t.Log(c1)
}