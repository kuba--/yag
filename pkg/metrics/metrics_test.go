package metrics

import (
	"encoding/json"
	"fmt"
)

func ExampleConsolidateByAvg() {
	const js = `["[1720.0, 1370846820]", "[1637.0, 1370846880]", "[1669.0, 1370846930]", "[1651.0, 1370847000]", "[1425.0, 1370847010]"]`

	var data []interface{}

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateByAvg(data, 1370846820, 1370847060, 60)
		for _, dp := range datapoints {
			fmt.Printf("[%.1f, %.0f]\n", dp[0], dp[1])
		}
	}

	// Output:
	//[1720.0, 1370846820]
	//[1653.0, 1370846880]
	//[0.0, 1370846940]
	//[1538.0, 1370847000]
	//[0.0, 1370847060]
}

func ExampleConsolidateBySum() {
	const js = `["[1, 0]", "[1, 1]", "[1, 2]", "[1, 3]", "[1, 4]", "[1, 5]", "[1, 6]", "[3, 20]", "[4, 30]", "[4, 34]"]`
	var data []interface{}

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateBySum(data, 0, 50, 10)
		for _, dp := range datapoints {
			fmt.Printf("[%.1f, %.0f]\n", dp[0], dp[1])
		}
	}

	// Output:
	//[7.0, 0]
	//[0.0, 10]
	//[3.0, 20]
	//[8.0, 30]
	//[0.0, 40]
	//[0.0, 50]
}

func ExampleConsolidateByMax() {
	js := `["[1.0, 0]", "[1.1, 1]", "[1.2, 2]", "[1.3, 3]", "[1.4, 4]", "[1, 5]", "[1.6, 6]", "[3, 20]", "[4.1, 30]", "[4, 34]"]`
	var data []interface{}

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateByMax(data, 0, 50, 10)
		for _, dp := range datapoints {
			fmt.Printf("[%.1f, %.0f]\n", dp[0], dp[1])
		}
	}

	// Output:
	//[1.6, 0]
	//[0.0, 10]
	//[3.0, 20]
	//[4.1, 30]
	//[0.0, 40]
	//[0.0, 50]
}

func ExampleConsolidateByMin() {
	js := `["[1.0, 0]", "[1.1, 1]", "[1.2, 2]", "[1.3, 3]", "[1.4, 4]", "[1, 5]", "[1.6, 6]", "[3, 20]", "[4.1, 30]", "[4, 34]"]`
	var data []interface{}

	if err := json.Unmarshal([]byte(js), &data); err != nil {
		fmt.Println(err)
	} else {
		datapoints := consolidateByMin(data, 0, 50, 10)
		for _, dp := range datapoints {
			fmt.Printf("[%.1f, %.0f]\n", dp[0], dp[1])
		}
	}

	// Output:
	//[1.0, 0]
	//[0.0, 10]
	//[3.0, 20]
	//[4.0, 30]
	//[0.0, 40]
	//[0.0, 50]
}
