package service

import (
	"fmt"
	"go-task-tracker/model"
	"log/slog"
	"time"
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

func (s *TaskService) AddTask(newTask model.CreateTask) error {

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

func (s *TaskService) GetTasks(status model.TaskStatus, description string) ([]model.Task, error) {

	tasks, err := s.repository.GetAllTasks()
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to get tasks using filters status=%d, description=%s", status, description), slog.Any("err", err))
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	if status == -1 && description == "" {
		return tasks, nil
	}

	tasksFiltered := make([]model.Task, 0)
	for _, task := range tasks {
		if (description == "" || description == task.Description) && (status == -1 || task.Status == status) {
			tasksFiltered = append(tasksFiltered, task)
		}
	}

	return tasksFiltered, nil
}

func (s *TaskService) UpdateTask(taskId int, taskToUpdate model.UpdateTask) error {
	s.log.Info(fmt.Sprintf("Updating task %d with values %+v", taskId, taskToUpdate))
	if err := s.repository.UpdateTask(taskId, taskToUpdate); err != nil {
		s.log.Error(fmt.Sprintf("error when updating task: %s", err))
		return fmt.Errorf("failed to update task: %w", err)
	}
	s.log.Info(fmt.Sprintf("Task %d updated with values %+v", taskId, taskToUpdate))
	return nil
}
