package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestDataRaceConditionsWithMutex(t *testing.T) {
	var state int32
	var mu sync.RWMutex

	for i := 0; i < 10; i++ {
		go func(i int) {
			mu.Lock()
			state += int32(i)
			// ... some more business logic
			mu.Unlock()
		}(i)
	}
}

func TestDataRaceConditionsWithAtomicValues(t *testing.T) {
	var state int32

	for i := 0; i < 10; i++ {
		go func(i int) {
			// separates memory for doing one operation at a time
			atomic.AddInt32(&state, int32(i))
			// ... some more business logic
		}(i)
	}
}

// run: go test race_test.go --race
// to check for race conditions

// atomic values will be faster than mutexes because they only lock
// the memory for the operation being done, while mutexes lock the
// memory for the entire block of code
