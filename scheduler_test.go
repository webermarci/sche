package sche

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	callCounter := 0

	scheduler := NewScheduler().
		WithDurationLimit(20 * time.Millisecond).
		Operation(func(task *Task, data interface{}) error {
			callCounter++
			return errors.New("next")
		}).
		Errors(func(err error) {
			t.Fatal(err)
		})

	scheduler.Schedule(t.Name(), 10*time.Millisecond)

	time.Sleep(30 * time.Millisecond)

	if callCounter != 2 {
		t.Fatalf("did not run 2 times: %d", callCounter)
	}
}

func TestSchedulerWithPersistence(t *testing.T) {
	defer func() {
		err := os.RemoveAll(t.Name())
		if err != nil {
			t.Fatal(err)
		}
	}()

	callCounter := 0

	scheduler := NewScheduler().
		WithDurationLimit(20 * time.Millisecond).
		WithPersistence(t.Name()).
		Operation(func(task *Task, data interface{}) error {
			callCounter++
			return errors.New("next")
		}).
		Errors(func(err error) {
			t.Fatal(err)
		})

	scheduler.Schedule(t.Name(), 10*time.Millisecond)

	time.Sleep(30 * time.Millisecond)

	if callCounter != 2 {
		t.Fatalf("did not run 2 times: %d", callCounter)
	}
}

func TestLoadingTasks(t *testing.T) {
	defer func() {
		err := os.RemoveAll(t.Name())
		if err != nil {
			t.Fatal(err)
		}
	}()

	scheduler := NewScheduler().
		WithDurationLimit(time.Second).
		WithPersistence(t.Name()).
		Operation(func(task *Task, data interface{}) error {
			return errors.New("next")
		}).
		Errors(func(err error) {
			t.Fatal(err)
		})

	task := newTask(t.Name(), 10*time.Millisecond, scheduler.durationLimit)
	scheduler.persistTask(task)

	called := false

	newScheduler := NewScheduler().
		WithDurationLimit(time.Second).
		WithPersistence(t.Name()).
		Operation(func(task *Task, data interface{}) error {
			called = true
			return nil
		}).
		Errors(func(err error) {
			t.Fatal(err)
		})

	if err := newScheduler.Load(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(20 * time.Millisecond)

	if !called {
		t.Fatal("task should have been called by now")
	}
}

func TestSchedulerWithoutOperation(t *testing.T) {
	hasError := false

	scheduler := NewScheduler().
		WithDurationLimit(20 * time.Millisecond).
		Errors(func(err error) {
			hasError = true
		})

	scheduler.Schedule(t.Name(), 10*time.Millisecond)

	time.Sleep(30 * time.Millisecond)

	if !hasError {
		t.Fatal("expected error")
	}
}
