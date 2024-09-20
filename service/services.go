package service

import (
	"encoding/json"
	"fmt"
	"go-task-tracker/model"
	"os"
)

func AddTask(task model.Task, path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var tasks []model.Task
	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &tasks)
		if err != nil {
			return fmt.Errorf("failed to unmarshal tasks: %w", err)
		}
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
