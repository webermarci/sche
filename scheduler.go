package sche

import (
	"errors"
	"time"

	"github.com/webermarci/pantry"
)

type Scheduler struct {
	operation     func(task *Task, data interface{}) error
	errors        func(err error)
	persisted     bool
	pantry        pantry.Pantry
	durationLimit time.Duration
}

func (scheduler *Scheduler) Operation(operation func(task *Task, data interface{}) error) *Scheduler {
	scheduler.operation = operation
	return scheduler
}

func (scheduler *Scheduler) Errors(errors func(err error)) *Scheduler {
	scheduler.errors = errors
	return scheduler
}

func (scheduler *Scheduler) WithPersistence(directory string) *Scheduler {
	scheduler.persisted = true
	scheduler.pantry = *pantry.New(&pantry.Options{
		CleaningInterval:     10 * time.Second,
		PersistenceDirectory: directory,
	}).Type(Task{})
	return scheduler
}

func (scheduler *Scheduler) WithDurationLimit(duration time.Duration) *Scheduler {
	scheduler.durationLimit = duration
	return scheduler
}

func (scheduler *Scheduler) Load() error {
	if err := scheduler.pantry.Load(); err != nil {
		return err
	}

	for _, value := range scheduler.pantry.GetAll() {
		task := value.(Task)
		task.run(scheduler)
	}

	return nil
}

func (scheduler *Scheduler) Schedule(data interface{}, periodicity time.Duration) {
	if scheduler.operation == nil && scheduler.errors != nil {
		scheduler.errors(errors.New("missing operation"))
		return
	}

	task := newTask(data, periodicity, scheduler.durationLimit)

	if scheduler.persisted {
		go func() {
			scheduler.pantry.Set(task.ID, task, scheduler.durationLimit).Persist()
		}()
	}

	go task.run(scheduler)
}

func (scheduler *Scheduler) persistTask(task *Task) error {
	if scheduler.persisted {
		return scheduler.pantry.Set(task.ID, *task, scheduler.durationLimit).Persist()
	}
	return nil
}

func (scheduler *Scheduler) removeTask(key string) error {
	if scheduler.persisted {
		return scheduler.pantry.Remove(key).Persist()
	}
	return nil
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}
