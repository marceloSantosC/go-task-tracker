package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-task-tracker/model"
	"io/fs"
	"os"
)

func AddTask(task model.Task, path string) error {

	tasks, err := GetAllTasks(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	tasks = append(tasks, task)
	tasksJson, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task %d: %w", task.Id, err)
	}

	err = os.WriteFile(path, tasksJson, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func UpdateTask(taskId int, description string, path string) error {
	tasks, err := GetAllTasks(path)
	if err != nil {
		return fmt.Errorf("failed to retrieve tasks: %w", err)
	}

	if len(tasks) == 0 {
		return errors.New("task with id does not exists")
	}

	taskExists := false
	for i := range tasks {
		if tasks[i].Id == taskId {
			tasks[i].Description = description
			taskExists = true
			break
		}
	}

	if !taskExists {
		return errors.New("task with id does not exists")
	}

	tasksJson, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task %d: %w", taskId, err)
	}

	err = os.WriteFile(path, tasksJson, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func GetAllTasks(path string) ([]model.Task, error) {
	fileData, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var tasks []model.Task
	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &tasks)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
		}
	}
	return tasks, nil
}

func DeleteTask(taskId int, path string) error {

	tasks, err := GetAllTasks(path)
	if err != nil {
		return fmt.Errorf("failed to delete task %d: %w", taskId, err)
	}

	if len(tasks) == 1 && tasks[0].Id == taskId {
		tasks = make([]model.Task, 0)
		if err = writeTaskList(tasks, path); err != nil {
			return fmt.Errorf("failed to delete task %d: %w", taskId, err)
		}
		return nil
	}

	taskToDeleteIndex := -1
	for i, task := range tasks {
		if task.Id == taskId {
			taskToDeleteIndex = i
			break
		}
	}

	if taskToDeleteIndex == -1 {
		return fmt.Errorf("failed to delete: task with id %d does not exists", taskId)
	}

	lastTaskIndex := len(tasks) - 1

	if taskToDeleteIndex != lastTaskIndex {
		tasks[taskToDeleteIndex] = tasks[lastTaskIndex]
		tasks = tasks[:lastTaskIndex]
	} else {
		tasks = tasks[:len(tasks)-1]
	}

	if err = writeTaskList(tasks, path); err != nil {
		return fmt.Errorf("failed to delete task %d: %w", taskId, err)
	}

	return nil
}

func writeTaskList(tasks []model.Task, path string) error {
	tasksJson, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to write tasks to file %s: %w", path, err)
	}

	err = os.WriteFile(path, tasksJson, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
