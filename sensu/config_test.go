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

	d1, _ = extend(d1, d2)
	t.Log(d1)
}