package metrics

import (
	"encoding/json"
	"fmt"
	"testing"
)

func ExampleSumSeries() {
	const data = `[
	{"target": "rest.200", "datapoints": [[1720.0, 1370846820], [1637.0, 1370846880], [1669.0, 1370846940], [1651.0, 1370847000], [1425.0, 1370847060]]},
	{"target": "rest.204", "datapoints": [[1.0, 1370846820],[1.0, 1370846940], [1.0, 1370847000]]},
	{"target": "rest.201", "datapoints": [[-1.0, 1370846820],[-1.0, 1370846940], [-1.0, 1370847000]]}]`

	const expected = `[{"target": "sum", 
	"datapoints": [[1720.00, 1370846820], [1637.00, 1370846880], [1669.00, 1370846940], [1651.00, 1370847000], [1425.00, 1370847060]]}]`

	var m []*Metrics
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		fmt.Println(err)
	} else {
		m := new(Api).sumSeries(m)[0]
		for _, dp := range m.Datapoints {
			val := "null"
			if dp[0] != nil {
				val = fmt.Sprintf("%.2f", *dp[0])
			}
			fmt.Printf("[%s, %.0f]\n", val, *dp[1])
		}
	}
	// Output:
	// [1720.00, 1370846820]
	// [1637.00, 1370846880]
	// [1669.00, 1370846940]
	// [1651.00, 1370847000]
	// [1425.00, 1370847060]
}

func TestEqual(t *testing.T) {

	const d1 = `["[1720.0, 1370846820]", "[1637.0, 1370846830]", "[1669.0, 1370846910]", "[1651.0, 1370847000]", "[1425.0, 1370847060]"]`
	const d2 = `["[1.0, 1370846820]","[1.0, 1370846910]", "[1.0, 1370847000]"]`
	const d3 = `["[-1.0, 1370846820]","[-1.0, 1370846910]", "[-1.0, 1370847000]"]`

	const d123 = `["[1720.00, 1370846820]", "[1637.00, 1370846830]", "[1669.00, 1370846910]", "[1651.00, 1370847000]", "[1425.00, 1370847060]"]`

	var isEq = func(m1, m2 *Metrics) bool {
		if m2 == nil {
			return false
		}

		d1, d2 := m1.Datapoints, m2.Datapoints
		ld1, ld2 := len(d1), len(d2)
		if ld1 != ld2 {
			return false
		}

		for i := 0; i < ld1; i++ {
			v1, t1 := d1[i][0], d1[i][1]
			v2, t2 := d2[i][0], d2[i][1]
			if *t1 != *t2 || (v1 != nil && v2 != nil && *v1 != *v2) || (v1 == nil && v2 != nil) || (v1 != nil && v2 == nil) {
				return false
			}
		}
		return true
	}

	var m1 []*Metrics
	var data1, data2, data3, data123 []interface{}

	if err := json.Unmarshal([]byte(d1), &data1); err != nil {
		fmt.Println(err)
	} else {
		m1 = append(m1, newMetrics("d1", "d1", consolidateBy(data1, 1370846820, 1370847060, 60, "avg")))
	}

	if err := json.Unmarshal([]byte(d2), &data2); err != nil {
		fmt.Println(err)
	} else {
		m1 = append(m1, newMetrics("d2", "d2", consolidateBy(data2, 1370846820, 1370847060, 60, "avg")))
	}

	if err := json.Unmarshal([]byte(d3), &data3); err != nil {
		fmt.Println(err)
	} else {
		m1 = append(m1, newMetrics("d3", "d3", consolidateBy(data3, 1370846820, 1370847060, 60, "avg")))
	}

	var m2 []*Metrics
	if err := json.Unmarshal([]byte(d123), &data123); err != nil {
		t.Error(err)
	} else {
		m2 = append(m2, newMetrics("d123", "d123", consolidateBy(data123, 1370846820, 1370847060, 60, "avg")))
	}

	m1 = new(Api).sumSeries(m1)
	if !isEq(m2[0], m1[0]) {
		t.Fatalf("\nexpected: %t\n\nwas: %t", m2[0], m1[0])
	}
}
