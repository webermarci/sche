package sche

import (
	"errors"
	"time"
)

type Scheduler[T any] struct {
	operation     func(data T)
	errors        func(err error)
	durationLimit time.Duration
}

func (scheduler *Scheduler[T]) Operation(operation func(data T)) *Scheduler[T] {
	scheduler.operation = operation
	return scheduler
}

func (scheduler *Scheduler[T]) Errors(errors func(err error)) *Scheduler[T] {
	scheduler.errors = errors
	return scheduler
}

func (scheduler *Scheduler[T]) WithDurationLimit(duration time.Duration) *Scheduler[T] {
	scheduler.durationLimit = duration
	return scheduler
}

func (scheduler *Scheduler[T]) Schedule(data T, periodicity time.Duration) {
	if scheduler.operation == nil && scheduler.errors != nil {
		scheduler.errors(errors.New("missing operation"))
		return
	}

	task := newTask(data, periodicity, scheduler.durationLimit)
	go task.run(scheduler)
}

func NewScheduler[T any]() *Scheduler[T] {
	return &Scheduler[T]{}
}
