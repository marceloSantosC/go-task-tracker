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

func removeFile(fileName string) {
	if err := os.Remove(fileName); err != nil {
		log.Panicf("Failed to remove file %s: %v", fileName, err)
	}
}

func Test_AddTask(t *testing.T) {
	const fileName = "test_task.json"
	defer removeFile(fileName)

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
	defer removeFile(fileName)

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
	defer removeFile(fileName)

	task, err := addTaskOrFail(fileName, t)
	if err != nil {
		return
	}
	newDescription := "New task description"

	if err = UpdateTask(task.Id, newDescription, fileName); err != nil {
		t.Errorf("expected to get no errors from CreateOrUpdateTask call but got error: %v", err)
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
	defer removeFile(fileName)

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
	defer removeFile(fileName)

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
	defer removeFile(fileName)

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

func Test_UpdateStatus(t *testing.T) {
	const fileName = "test_task.json"
	defer removeFile(fileName)

	task, err := addTaskOrFail(fileName, t)
	if err != nil {
		return
	}

	var testTable = []struct {
		taskId         int
		expectedStatus model.TaskStatus
		path           string
	}{
		{task.Id, model.TODO, fileName},
		{task.Id, model.InProgress, fileName},
		{task.Id, model.Done, fileName},
	}

	for _, d := range testTable {

		testName := fmt.Sprintf("for input taskId=%d, status=%d(%s) and fileName=%s expect status to be %d",
			d.taskId, d.expectedStatus, d.expectedStatus.String(), d.path, d.expectedStatus)

		t.Run(testName, func(t *testing.T) {
			if err = UpdateTaskStatus(d.taskId, d.expectedStatus, d.path); err != nil {
				t.Errorf("expected no error from UpdateTaskStatus call but got \"%s\"", err)
				return
			}

			retrievedTasks, err := getTasksFromFileOrFail(fileName, t)
			if err != nil {
				return
			}

			var retrievedTask model.Task
			for _, retrievedTask = range retrievedTasks {
				if retrievedTask.Id == d.taskId {
					break
				}
			}

			if retrievedTask == (model.Task{}) || retrievedTask.Id != d.taskId {
				t.Errorf("expected to retrieve task with id %d, but got no tasks", d.taskId)
				return
			}

			if retrievedTask.Status != d.expectedStatus {
				t.Errorf("expected task %d to have status %d(%s), but got %d(%s)", d.taskId, d.expectedStatus,
					d.expectedStatus.String(), retrievedTask.Status, retrievedTask.Status.String())
				return
			}

		})

	}
}

func Test_GetAllTasksByStatusEqualsTrue(t *testing.T) {
	const fileName = "test_task.json"
	defer removeFile(fileName)

	tasks := make([]model.Task, 3)
	for i := range 3 {
		task := model.Task{
			Id:          i,
			Description: fmt.Sprintf("Test %d", i),
			Status:      model.TaskStatus(i),
			CreatedAt:   model.DateTime(time.Now()),
			UpdatedAt:   model.DateTime(time.Now()),
		}
		tasks[i] = task

		if err := AddTask(task, fileName); err != nil {
			t.Errorf("expected to get no errors from AddTask call but got error: %v", err)
			return
		}
	}

	var testData = []struct{ want model.TaskStatus }{
		{want: model.TODO}, {want: model.InProgress}, {want: model.Done},
	}

	for _, d := range testData {

		var testName = fmt.Sprintf("should get all tasks with status %s given inputs status=%s, equals=%t, path=%s",
			d.want.String(), d.want.String(), true, fileName)

		t.Run(testName, func(t *testing.T) {
			returnedTasks, err := GetAllTasksByStatus(d.want, true, fileName)
			if err != nil {
				t.Errorf("expected no errors from GetAllTasksByStatus call got \"%s\"", err)
				return
			}

			if nOfTasks := len(returnedTasks); nOfTasks != 1 {
				t.Errorf("expected one task returned from GetAllTasksByStatus got %d tasks", nOfTasks)
				return
			}

			if status := returnedTasks[0].Status; status != d.want {
				t.Errorf("expected task to have status %d(%s) got %d(%s)", status, status.String(), d.want, d.want.String())
				return
			}

		})
	}
}

func Test_GetAllTasksByStatusEqualsFalse(t *testing.T) {
	const fileName = "test_task.json"
	defer removeFile(fileName)

	tasks := make([]model.Task, 3)
	for i := range 3 {
		task := model.Task{
			Id:          i,
			Description: fmt.Sprintf("Test %d", i),
			Status:      model.TaskStatus(i),
			CreatedAt:   model.DateTime(time.Now()),
			UpdatedAt:   model.DateTime(time.Now()),
		}
		tasks[i] = task

		if err := AddTask(task, fileName); err != nil {
			t.Errorf("expected to get no errors from AddTask call but got error: %v", err)
			return
		}
	}

	var testData = []struct{ doNotWant model.TaskStatus }{
		{doNotWant: model.TODO}, {doNotWant: model.InProgress}, {doNotWant: model.Done},
	}

	for _, d := range testData {

		var testName = fmt.Sprintf("should get all tasks with status NOT EQUALS %s given inputs status=%s, equals=%t, path=%s",
			d.doNotWant.String(), d.doNotWant.String(), false, fileName)

		t.Run(testName, func(t *testing.T) {
			returnedTasks, err := GetAllTasksByStatus(d.doNotWant, false, fileName)
			if err != nil {
				t.Errorf("expected no errors from GetAllTasksByStatus call got \"%s\"", err)
				return
			}

			if nOfTasks := len(returnedTasks); nOfTasks != 2 {
				t.Errorf("expected two tasks returned from GetAllTasksByStatus got %d tasks", nOfTasks)
				return
			}

			for _, returnedTask := range returnedTasks {
				if returnedTask.Status == d.doNotWant {
					t.Errorf("expected task to not have status %d(%s) got %+v", d.doNotWant, d.doNotWant.String(), returnedTask)
					return
				}
			}

		})
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
