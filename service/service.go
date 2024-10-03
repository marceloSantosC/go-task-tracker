package service

import (
	"fmt"
	"go-task-tracker/model"
	"log/slog"
)

type TaskRepository interface {
	AddTask(task model.Task) error

	UpdateTask(taskId int, updatedTask model.UpdateTask) error

	GetAllTasks() ([]model.Task, error)

	DeleteTask(taskId int) error
}

type Error struct {
	UserMsg string
	err     error
}

func (e Error) Error() string {
	return e.err.Error()
}

func NewError(err error, userMsg string) Error {
	return Error{
		UserMsg: userMsg,
		err:     err,
	}
}

type TaskService struct {
	repository TaskRepository
	log        slog.Logger
}

func NewTaskService(repository TaskRepository, log *slog.Logger) TaskService {
	return TaskService{repository: repository, log: *log}
}

func (s *TaskService) AddTask(task model.Task) error {
	err := s.repository.AddTask(task)
	if err != nil {
		err = fmt.Errorf("failed to create task: %w", err)
		s.log.Error(err.Error())
		return NewError(err, "error when creating user, check logs for more info")
	}
	return nil
}
