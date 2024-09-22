package repository

import (
	"encoding/json"
	"fmt"
	"go-task-tracker/model"
	"os"
	"testing"
	"time"
)

func Test_NewTaskRepositoryFile_WithCreatedFile(t *testing.T) {
	const fileName = "NewTaskRepositoryFile.json"
	defer removeTestFile(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("failed to create test file %s error: \"%s\"", fileName, err)
	}
	defer file.Close()

	if _, err = file.WriteString("[\n\n]"); err != nil {
		t.Fatalf("failed to write to test file %s error: \"%s\"", fileName, err)
	}

	repository, err := NewTaskRepositoryFile(fileName)
	if err != nil {
		t.Fatalf("expected call to NewTaskRepositoryFile to return no errors, but got \"%v\"", err)
	}

	if repository.path != fileName {
		t.Errorf("expect path to be %s but was %s", fileName, repository.path)
	}

	if repository.offset == 0 {
		t.Error("expect offset to not be 0")
	}
}

func Test_NewTaskRepositoryFile_WithNoFile(t *testing.T) {
	const fileName = "NewTaskRepositoryFile.json"
	defer removeTestFile(fileName)

	repository, err := NewTaskRepositoryFile(fileName)
	if err != nil {
		t.Fatalf("expected call to NewTaskRepositoryFile to return no errors, but got \"%v\"", err)
	}

	if repository.path != fileName {
		t.Errorf("expect path to be %s but was %s", fileName, repository.path)
	}

	if repository.offset != 2 {
		t.Errorf("expect offset to be 1 but was %d", repository.offset)
	}
}

func Test_AddTask(t *testing.T) {
	const fileName = "AddTask.json"
	defer removeTestFile(fileName)

	r, err := NewTaskRepositoryFile(fileName)
	if err != nil {
		t.Fatalf("failed to create TaskRepositoryFile")
	}
	task := model.Task{
		Id:          1,
		Description: "Test",
		Status:      0,
		CreatedAt:   model.DateTime(time.Now()),
		UpdatedAt:   model.DateTime(time.Now()),
	}

	if err = r.AddTask(task); err != nil {
		t.Fatalf("failed to call AddTask: \"%v\"", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("failed to create TaskRepositoryFile")
	}
	defer file.Close()

	var tasks []model.Task
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&tasks); err != nil {
		t.Fatalf("failed to parse json from file %s: \"%v\"", fileName, err)
	}

	if nOfTasks := len(tasks); nOfTasks != 1 {
		t.Errorf("expected one task to be created, got %d", nOfTasks)
	}

	if task = tasks[0]; task.Id != r.sequenceId {
		t.Errorf("expected task id to be %d got %d", r.sequenceId, task.Id)
	}
}

func removeTestFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		panic(fmt.Errorf("failed to remove file %s: %w", fileName, err))
	}
}
