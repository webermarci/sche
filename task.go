package sche

import (
	"strconv"
	"time"
)

type Task[T any] struct {
	ID          string
	Data        T
	Periodicity time.Duration
	TimeLimit   time.Time
}

func (task *Task[T]) run(scheduler *Scheduler[T]) {
	if time.Now().After(task.TimeLimit) {
		return
	}

	scheduler.operation(task.Data)
	time.Sleep(task.Periodicity)
	task.run(scheduler)
}

func newTask[T any](data T, periodicity time.Duration, durationLimit time.Duration) *Task[T] {
	now := time.Now()
	return &Task[T]{
		ID:          strconv.Itoa(int(now.UnixNano())),
		Data:        data,
		Periodicity: periodicity,
		TimeLimit:   now.Add(durationLimit),
	}
}
