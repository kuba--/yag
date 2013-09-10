package funcexp

import (
	"log"
	"strconv"
	"testing"
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

func TestEval_Val(t *testing.T) {

	const expr = "123456789"
	const expected = 123456789

	api := new(IntApi)
	received, _ := Eval(expr, 0, 0, api).(int)
	if received != expected {
		t.Errorf("TestEval(%s) failed.\tExpected:%d\tReceived:%d\n", expr, expected, received)
	}
}

func TestEval_Sum(t *testing.T) {

	const expr = "sum(sum(1, 2), sub(2, 1), -13)"
	const expected = -9

	api := new(IntApi)
	received, _ := Eval(expr, 0, 0, api).(int)
	if received != expected {
		t.Errorf("TestEval(%s) failed.\tExpected:%d\tReceived:%d\n", expr, expected, received)
	}
}

func TestEval_SumOfSums(t *testing.T) {

	const expr = "sum(sum(sum(sum(sum(sum(1,2),-3),4), 5),-9),0)"
	const expected = 0

	api := new(IntApi)
	received, _ := Eval(expr, 0, 0, api).(int)
	if received != expected {
		t.Errorf("TestEval(%s) failed.\tExpected:%d\tReceived:%d\n", expr, expected, received)
	}
}
