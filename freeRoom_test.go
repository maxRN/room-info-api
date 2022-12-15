package main

import (
	"testing"
	"time"
)

var LOCATION = time.Local
var SECOND_DS = time.Date(2022, 12, 15, 9, 50, 0, 0, LOCATION)

func TestSecondDs(t *testing.T) {
	if getCurrentDs(SECOND_DS) != 2 {
		t.Fail()
	}
}
