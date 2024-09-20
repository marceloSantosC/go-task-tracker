package service

import (
	"encoding/json"
	"go-task-tracker/model"
	"log"
	"os"
	"testing"
	"time"
)

func Test_AddTask(t *testing.T) {

	fileName := "test_task.json"
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

	file, err := os.Open(fileName)
	if err != nil {
		t.Errorf("expected to open file %s but got error: %v", fileName, err)
		return
	}
	defer func() {
		_ = file.Close()
		if err = os.Remove(fileName); err != nil {
			log.Printf("Failed to remove file %s: %v", fileName, err)
		}
	}()

	var tasks []model.Task
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	for decoder.More() {
		err = decoder.Decode(&tasks)
		if err != nil {
			t.Errorf("expected to decode json from file but got error: %v", err)
			return
		}
	}

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
	fileName := "test_task.json"

	if err := addTaskOrFail(fileName, t); err != nil {
		return
	}

	tasks, err := GetAllTasks(fileName)
	if err != nil {
		t.Errorf("expected to get no errors from GetAllTasks call but got error: %v", err)
		return
	}

	if len(tasks) < 1 {
		t.Errorf("expected to retrieve one task, but got: %d", len(tasks))
		return
	}

}

/*
func Test_UpdateTask(t *testing.T) {

	fileName := "test_task.json"

	if err := addTaskOrFail(fileName, t); err != nil {
		return
	}

	newDescription := "New task description"

}
*/

func addTaskOrFail(fileName string, t *testing.T) error {
	task := model.Task{
		Id:          1,
		Description: "Test",
		Status:      0,
		CreatedAt:   model.DateTime(time.Now()),
		UpdatedAt:   model.DateTime(time.Now()),
	}
	if err := AddTask(task, fileName); err != nil {
		t.Errorf("expected to get no errors from AddTask call but got error: %v", err)
		return err
	}
	return nil
}
