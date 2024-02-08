package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	eatTime = 0 * time.Second   // time duration to eat for 1 philosopher
	thinkTime = 0 * time.Second // time duration to think for 1 philosopher
	sleepTime = 0 * time.Second // time to sleep

	for i := 0; i < 10; i++ {
		orderFinished = []string{}
		dine()
		if len(orderFinished) != 5 {
			t.Errorf("incorrect length of a slice, expected 5 but got %d\n", len(orderFinished))
		}
	}

}

func Test_dineWithVaryingDelays(t *testing.T) {

	var theTests = []struct {
		name  string
		delay time.Duration
	}{
		{"zero_delay", time.Second * 0},
		{"quarter_second_delay", time.Millisecond * 250},
		{"half_second_delay", time.Millisecond * 500},
	}

	for _, e := range theTests {
		orderFinished = []string{}
		eatTime = e.delay
		sleepTime = e.delay
		thinkTime = e.delay

		dine()
		if len(orderFinished) != 5 {
			t.Errorf("%s: incorrect length of a slice, expected 5 but got %d\n", e.name, len(orderFinished))
		}
	}

}
