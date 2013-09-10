package metrics

import (
	"encoding/json"
	"log"
	"testing"
)

func TestApiSum(t *testing.T) {
	const data = `[
	{"target": "rest.200", "datapoints": ["[1720.0, 1370846820]", "[1637.0, 1370846880]", "[1669.0, 1370846940]", "[1651.0, 1370847000]", "[1425.0, 1370847060]"]},
	{"target": "rest.204", "datapoints": ["[1.0, 1370846820]","[1.0, 1370846940]", "[1.0, 1370847000]"]},
	{"target": "rest.201", "datapoints": ["[-1.0, 1370846820]","[-1.0, 1370846940]", "[-1.0, 1370847000]"]}]`

	const expected = `[{"target": "sum", 
	"datapoints": ["[1720.00, 1370846820]", "[1637.00, 1370846880]", "[1669.00, 1370846940]", "[1651.00, 1370847000]", "[1425.00, 1370847060]"]}]`

	var dm []map[string]interface{}
	err := json.Unmarshal([]byte(data), &dm)
	if err != nil {
		t.Error(err)
	}
	var em []map[string]interface{}
	err = json.Unmarshal([]byte(expected), &em)
	if err != nil {
		t.Error(err)
	}

	api := NewApi(-1)
	was := api.sum(newMetrics("data", dm))
	exp := newMetrics("expected", em)

	if !exp[0].isEqual(was[0]) {
		log.Fatalf("\nexpected: %t\n\nwas: %t", exp[0], was[0])
	}
}
