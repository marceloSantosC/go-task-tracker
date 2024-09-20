package service

import (
	"encoding/json"
	"go-task-tracker/model"
	"log"
	"os"
	"testing"
	"time"
)

var fileName = "test_task.json"

func TestMain(m *testing.M) {
	code := m.Run()
	shutdown(fileName)
	os.Exit(code)
}

func shutdown(fileName string) {
	if err := os.Remove(fileName); err != nil {
		log.Printf("Failed to remove file %s: %v", fileName, err)
	}
}

func Test_AddTask(t *testing.T) {

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

	tasks, err := getTasksFromFileOrFail(t)

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
	task, err := addTaskOrFail(fileName, t)
	if err != nil {
		return
	}
	newDescription := "New task description"

	if err = UpdateTask(task.Id, newDescription, fileName); err != nil {
		t.Errorf("expected to get no errors from UpdateTask call but got error: %v", err)
		return
	}

	tasks, err := getTasksFromFileOrFail(t)
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

func getTasksFromFileOrFail(t *testing.T) ([]model.Task, error) {
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
