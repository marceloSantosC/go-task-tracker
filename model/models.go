package model

import (
	"fmt"
	"strings"
	"time"
)

type TaskStatus int

const (
	TODO TaskStatus = iota
	InProgress
	Done
)

func (t TaskStatus) String() string {
	return [...]string{"To do", "In progress", "Done"}[t]
}

func (t TaskStatus) EnumIndex() int {
	return int(t)
}

type DateTime time.Time

func (t *DateTime) String() string {
	return fmt.Sprintf(`"%s"`, time.Time(*t).Format("2006-01-02 15:04:05"))
}

func (t *DateTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(*t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t *DateTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse("2006-01-02 15:04:05", strings.Replace(string(b), `"`, "", 2))
	if err != nil {
		return err
	}
	*t = DateTime(date)
	return
}

type Task struct {
	Id          int        `json:"Id"`
	Description string     `json:"Description"`
	Status      TaskStatus `json:"Status"`
	CreatedAt   DateTime   `json:"CreatedAt"`
	UpdatedAt   DateTime   `json:"UpdatedAt"`
}

type UpdateTask struct {
	Description string     `json:"Description"`
	Status      TaskStatus `json:"Status"`
}
