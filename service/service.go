package service

import (
	"fmt"
	"go-task-tracker/model"
	"log/slog"
	"time"
)

type TaskRepository interface {
	AddTask(task model.Task) error

	UpdateTask(taskId int, updatedTask model.CreateOrUpdateTask) error

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

func (s *TaskService) AddTask(newTask model.CreateOrUpdateTask) error {

	task := model.Task{
		Description: newTask.Description,
		Status:      newTask.Status,
		CreatedAt:   model.DateTime(time.Now()),
		UpdatedAt:   model.DateTime(time.Now()),
	}

	err := s.repository.AddTask(task)
	if err != nil {
		err = fmt.Errorf("failed to create task: %w", err)
		s.log.Error(err.Error())
		return NewError(err, "error when creating user")
	}
	return nil
}
