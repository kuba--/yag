package api

import (
	"fmt"
	"log"
	"strconv"
)

type IntApi struct{}

func (api *IntApi) Value(name string, from int64, to int64) interface{} {
	if i, err := strconv.Atoi(name); err == nil {
		return i
	} else {
		log.Fatal(err)
	}
	return nil
}

func (api *IntApi) Call(name string, argv interface{}) interface{} {
	if ints, ok := argv.([]interface{}); ok {
		s := 0
		switch name {
		case "sum":
			for _, n := range ints {
				if ii, oo := n.(int); oo {
					s += ii
				}
			}
		case "sub":
			if n := len(ints); n > 0 {
				s, _ = ints[0].(int)
				for i := 1; i < n; i++ {
					if nn, oo := ints[i].(int); oo {
						s -= nn
					}
				}
			}
		}
		return s
	}

	return nil
}

func ExampleVal() {
	api := new(IntApi)
	received, _ := Eval("123456789", 0, 0, api).(int)
	fmt.Println(received)
	// Output:
	// 123456789
}

func ExampleSum() {
	api := new(IntApi)
	received, _ := Eval("sum(sum(1, 2), sub(2, 1), -13)", 0, 0, api).(int)
	fmt.Println(received)
	// Output:
	// -9
}

func ExampleSumOfSums() {
	api := new(IntApi)
	received, _ := Eval("sum(sum(sum(sum(sum(sum(1,2),-3),4), 5),-9),0)", 0, 0, api).(int)
	fmt.Println(received)
	// Output:
	// 0
}
