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
