package sche

import (
	"strconv"
	"time"
)

type Task struct {
	ID          string
	Data        interface{}
	Periodicity time.Duration
	TimeLimit   time.Time
}

func (task *Task) run(scheduler *Scheduler) error {
	if time.Now().After(task.TimeLimit) {
		scheduler.removeTask(task.ID)
		return nil
	}

	err := scheduler.operation(task, task.Data)
	if err != nil {
		scheduler.persistTask(task)
		time.Sleep(task.Periodicity)
		return task.run(scheduler)
	}

	scheduler.removeTask(task.ID)
	return nil
}

func newTask(data interface{}, periodicity time.Duration, durationLimit time.Duration) *Task {
	now := time.Now()
	return &Task{
		ID:          strconv.Itoa(int(now.UnixNano())),
		Data:        data,
		Periodicity: periodicity,
		TimeLimit:   now.Add(durationLimit),
	}
}
