package model

import "time"

type TaskStatus int

const (
	TODO TaskStatus = iota
	IN_PROGRESS
	DONE
)

func (t TaskStatus) String() string {
	return [...]string{"To do", "In progress", "Done"}[t]
}

func (t TaskStatus) EnumIndex() int {
	return int(t)
}

type Task struct {
	Id          string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
