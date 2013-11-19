package metrics

import (
	"encoding/json"
	"fmt"
)

func ExampleMaxDataPoints() {
	const js = `["[0,1384613389]","[0,1384613399]","[0,1384613409]","[0.5,1384613419]","[0.75,1384614209]"]`
	var (
		data          []interface{}
		maxDataPoints int   = 7
		from, to      int64 = 1384613389, 1384614209
	)
	step := consolidationStep(from, to, 60, maxDataPoints)

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateBy(data, from, to, step, "avg")
		for _, dp := range datapoints {
			val := "null"
			if dp[0] != nil {
				val = fmt.Sprintf("%.2f", *dp[0])
			}
			fmt.Printf("[%s, %.0f]\n", val, *dp[1])
		}
	}
	// Output:
	//[0.12, 1384613389]
	//[null, 1384613509]
	//[null, 1384613629]
	//[null, 1384613749]
	//[null, 1384613869]
	//[null, 1384613989]
	//[0.75, 1384614109]

}

func ExampleConsolidateByAvg() {
	const js = `["[1720.0, 1370846820]", "[1637.0, 1370846880]", "[1669.0, 1370846930]", "[1651.0, 1370847000]", "[1425.0, 1370847010]"]`

	var data []interface{}

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateBy(data, 1370846820, 1370847060, 60, "avg")
		for _, dp := range datapoints {
			val := "null"
			if dp[0] != nil {
				val = fmt.Sprintf("%.2f", *dp[0])
			}
			fmt.Printf("[%s, %.0f]\n", val, *dp[1])
		}
	}

	// Output:
	//[1720.00, 1370846820]
	//[1653.00, 1370846880]
	//[null, 1370846940]
	//[1538.00, 1370847000]
	//[null, 1370847060]
}

func ExampleConsolidateBySum() {
	const js = `["[1, 0]", "[1, 1]", "[1, 2]", "[1, 3]", "[1, 4]", "[1, 5]", "[1, 6]", "[3, 20]", "[4, 30]", "[4, 34]"]`
	var data []interface{}

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateBy(data, 0, 50, 10, "sum")
		for _, dp := range datapoints {
			val := "null"
			if dp[0] != nil {
				val = fmt.Sprintf("%.1f", *dp[0])
			}
			fmt.Printf("[%s, %.0f]\n", val, *dp[1])
		}
	}

	// Output:
	//[7.0, 0]
	//[null, 10]
	//[3.0, 20]
	//[8.0, 30]
	//[null, 40]
	//[null, 50]
}

func ExampleConsolidateByMax() {
	js := `["[1.0, 0]", "[1.1, 1]", "[1.2, 2]", "[1.3, 3]", "[1.4, 4]", "[1, 5]", "[1.6, 6]", "[3, 20]", "[4.1, 30]", "[4, 34]"]`
	var data []interface{}

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateBy(data, 0, 50, 10, "max")
		for _, dp := range datapoints {
			val := "null"
			if dp[0] != nil {
				val = fmt.Sprintf("%.1f", *dp[0])
			}
			fmt.Printf("[%s, %.0f]\n", val, *dp[1])
		}
	}

	// Output:
	//[1.6, 0]
	//[null, 10]
	//[3.0, 20]
	//[4.1, 30]
	//[null, 40]
	//[null, 50]
}

func ExampleConsolidateByMin() {
	js := `["[1.0, 0]", "[1.1, 1]", "[1.2, 2]", "[1.3, 3]", "[1.4, 4]", "[1, 5]", "[1.6, 6]", "[3, 20]", "[4.1, 30]", "[4, 34]"]`
	var data []interface{}

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateBy(data, 0, 50, 10, "min")
		for _, dp := range datapoints {
			val := "null"
			if dp[0] != nil {
				val = fmt.Sprintf("%.1f", *dp[0])
			}
			fmt.Printf("[%s, %.0f]\n", val, *dp[1])
		}
	}

	// Output:
	//[1.0, 0]
	//[null, 10]
	//[3.0, 20]
	//[4.0, 30]
	//[null, 40]
	//[null, 50]
}
