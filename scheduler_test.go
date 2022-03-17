package sche

import (
	"sync"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	type TestData struct {
		Text   string
		Number int
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	scheduler := NewScheduler[TestData]().
		WithDurationLimit(20 * time.Millisecond).
		Operation(func(data TestData) {
			wg.Done()
		}).
		Errors(func(err error) {
			t.Fatal(err)
		})

	scheduler.Schedule(TestData{
		Text:   t.Name(),
		Number: 42,
	}, 10*time.Millisecond)

	go func() {
		time.Sleep(50 * time.Millisecond)
		t.Fail()
	}()

	wg.Wait()
}

func TestSchedulerWithoutOperation(t *testing.T) {
	hasError := false

	scheduler := NewScheduler[string]().
		WithDurationLimit(20 * time.Millisecond).
		Errors(func(err error) {
			hasError = true
		})

	scheduler.Schedule(t.Name(), 10*time.Millisecond)

	time.Sleep(25 * time.Millisecond)

	if !hasError {
		t.Fatal("expected error")
	}
}
