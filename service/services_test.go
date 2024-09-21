package service

import (
	"encoding/json"
	"fmt"
	"go-task-tracker/model"
	"log"
	"os"
	"testing"
	"time"
)

func shutdown(fileName string) {
	if err := os.Remove(fileName); err != nil {
		log.Panicf("Failed to remove file %s: %v", fileName, err)
	}
}

func Test_AddTask(t *testing.T) {
	const fileName = "test_task.json"
	defer shutdown(fileName)

	task := model.Task{
		Id:          1,
		Description: "Test",
		Status:      0,
		CreatedAt:   model.DateTime(time.Now()),
		UpdatedAt:   model.DateTime(time.Now()),
	}

	err := AddTask(task, fileName)
	if err != nil {
		t.Errorf("expected to get no errors from AddTask call but got error: %v", err)
		return
	}

	tasks, err := getTasksFromFileOrFail(fileName, t)

	if tasksLen := len(tasks); tasksLen != 1 {
		t.Errorf("expected to retrieve one task, but got %d", tasksLen)
		return
	}

	if tasks[0].Id != task.Id {
		t.Errorf("expected to retrieve task the same task passed to AddNew but got %+v", tasks[0])
		return
	}

}

func Test_GetAllTasks(t *testing.T) {
	const fileName = "test_task.json"
	defer shutdown(fileName)

	if _, err := addTaskOrFail(fileName, t); err != nil {
		return
	}

	tasks, err := GetAllTasks(fileName)
	if err != nil {
		t.Errorf("expected to get no errors from GetAllTasks call but got error: %v", err)
		return
	}

	if len(tasks) != 1 {
		t.Errorf("expected to retrieve one task, but got: %d", len(tasks))
		return
	}

}

func Test_UpdateTask(t *testing.T) {
	const fileName = "test_task.json"
	defer shutdown(fileName)

	task, err := addTaskOrFail(fileName, t)
	if err != nil {
		return
	}
	newDescription := "New task description"

	if err = UpdateTask(task.Id, newDescription, fileName); err != nil {
		t.Errorf("expected to get no errors from UpdateTask call but got error: %v", err)
		return
	}

	tasks, err := getTasksFromFileOrFail(fileName, t)
	if err != nil {
		return
	}

	if len(tasks) == 0 {
		t.Errorf("expected to retrieve array with one task, but got array with %d tasks", len(tasks))
		return
	}

	var actualTask model.Task
	for _, actualTask = range tasks {
		if actualTask.Id == task.Id {
			break
		}
	}

	if actualTask == (model.Task{}) {
		t.Errorf("expected to retrieve task with id %d but got no tasks", task.Id)
		return
	}

	if actualTask.Description != newDescription {
		t.Errorf(`expected to retrieve task with description "%s" but got task with description "%s"`, newDescription, tasks[0].Description)
		return
	}

}

func Test_DeleteTask(t *testing.T) {
	const fileName = "test_task.json"
	defer shutdown(fileName)

	task, err := addTaskOrFail(fileName, t)
	if err != nil {
		return
	}

	taskId := task.Id

	err = DeleteTask(taskId, fileName)
	if err != nil {
		t.Errorf("expected no errors from DeleteTask call got: \"%s\"", err)
		return
	}

	tasks, err := getTasksFromFileOrFail(fileName, t)
	if err != nil {
		return
	}

	for _, returnedTask := range tasks {
		if returnedTask.Id == taskId {
			t.Errorf("expect to retrieve no tasks with id %d, got %+v", taskId, returnedTask)
			return
		}
	}
}

func Test_DeleteLastTaskWithMultipleTasksWritten(t *testing.T) {
	const fileName = "test_task.json"
	defer shutdown(fileName)

	createdTasks, err := addNTasksOrFail(fileName, 5, t)
	if err != nil {
		return
	}

	taskId := createdTasks[len(createdTasks)-1].Id

	err = DeleteTask(taskId, fileName)
	if err != nil {
		t.Errorf("expected no errors from DeleteTask call got: \"%s\"", err)
		return
	}

	retrievedTasks, err := getTasksFromFileOrFail(fileName, t)
	if err != nil {
		return
	}

	for _, retrievedTask := range retrievedTasks {
		if retrievedTask.Id == taskId {
			t.Errorf("expect to retrieve no tasks with id %d, got %+v", taskId, retrievedTask)
			return
		}
	}
}

func Test_DeleteTaskWithMultipleTasksWritten(t *testing.T) {
	const fileName = "test_task.json"
	defer shutdown(fileName)

	createdTasks, err := addNTasksOrFail(fileName, 5, t)
	if err != nil {
		return
	}

	taskId := createdTasks[2].Id

	err = DeleteTask(taskId, fileName)
	if err != nil {
		t.Errorf("expected no errors from DeleteTask call got: \"%s\"", err)
		return
	}

	retrievedTasks, err := getTasksFromFileOrFail(fileName, t)
	if err != nil {
		return
	}

	for _, retrievedTask := range retrievedTasks {
		if retrievedTask.Id == taskId {
			t.Errorf("expect to retrieve no tasks with id %d, got %+v", taskId, retrievedTask)
			return
		}
	}
}

func addNTasksOrFail(fileName string, numberOfTasks int, t *testing.T) ([]model.Task, error) {

	tasks := make([]model.Task, numberOfTasks)

	for i := range numberOfTasks {
		task := model.Task{
			Id:          i,
			Description: fmt.Sprintf("Test %d", i),
			Status:      0,
			CreatedAt:   model.DateTime(time.Now()),
			UpdatedAt:   model.DateTime(time.Now()),
		}

		tasks[i] = task

		if err := AddTask(task, fileName); err != nil {
			t.Errorf("expected to get no errors from AddTask call but got error: %v", err)
			return nil, err
		}
	}

	return tasks, nil
}

func addTaskOrFail(fileName string, t *testing.T) (model.Task, error) {
	task := model.Task{
		Id:          1,
		Description: "Test",
		Status:      0,
		CreatedAt:   model.DateTime(time.Now()),
		UpdatedAt:   model.DateTime(time.Now()),
	}
	if err := AddTask(task, fileName); err != nil {
		t.Errorf("expected to get no errors from AddTask call but got error: %v", err)
		return model.Task{}, err
	}
	return task, nil
}

func getTasksFromFileOrFail(fileName string, t *testing.T) ([]model.Task, error) {
	file, err := os.Open(fileName)
	if err != nil {
		t.Errorf("expected to open file %s but got error: %v", fileName, err)
		return nil, err
	}
	defer file.Close()

	var tasks []model.Task
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	for decoder.More() {
		err = decoder.Decode(&tasks)
		if err != nil {
			t.Errorf("expected to decode json from file but got error: %v", err)
			return nil, err
		}
	}

	return tasks, nil
}
