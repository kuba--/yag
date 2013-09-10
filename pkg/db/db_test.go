package db

import (
	"testing"
	"time"
)

func TestDbConn(t *testing.T) {
	conn := func(name string) {
		t.Log("\t", name)

		if c1, err1 := Client(); err1 != nil {
			t.Error(err1)
		} else {
			defer Release(c1)
		}
		time.Sleep(2 * time.Second)
	}

	go conn("c1")
	conn("c2")
	conn("c3")
}
