package repository

import (
	"bufio"
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

func Test_UpdateTask(t *testing.T) {
	const fileName = "Test_UpdateTask.json"
	defer removeTestFile(fileName)

	repository, err := NewTaskRepositoryFile(fileName)
	if err != nil {
		t.Fatalf("failed to create TaskRepositoryFile")
	}
	tasks := newTasks(3)
	addTasksToFileOrFail(tasks, fileName, t)

	description := "Lorem"
	status := model.TaskStatus(1)
	updatedTask := model.UpdateTask{
		Description: &description,
		Status:      &status,
	}

	task := tasks[0]

	if err = repository.UpdateTask(task.Id, updatedTask); err != nil {
		t.Fatalf("expected call to CreateTask to not return error, got \"%v\"", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("failed to open file %s", fileName)
	}
	defer file.Close()

	var tasksInFile []model.Task
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&tasksInFile); err != nil {
		t.Fatalf("failed to parse json from file %s: \"%v\"", fileName, err)
	}

	var taskInFile model.Task
	for _, taskInFile = range tasksInFile {
		if taskInFile.Id == task.Id {
			break
		}
	}

	if taskInFile == (model.Task{}) {
		t.Fatalf("expected to find task with id %d", task.Id)
	}

	if taskInFile.Description != *updatedTask.Description {
		t.Fatalf("expected task description to be %s, got %s", *updatedTask.Description, taskInFile.Description)
	}

	if taskInFile.UpdatedAt == task.UpdatedAt {
		t.Fatalf("expected task updated at field to change, but was the same as before")
	}

}

func Test_GetAllTasks(t *testing.T) {

	fileName := "Test_GetAllTasks.json"
	defer removeTestFile(fileName)

	tasks := newTasks(5)
	addTasksToFileOrFail(tasks, fileName, t)

	repository, err := NewTaskRepositoryFile(fileName)
	if err != nil {
		t.Fatalf("failed to create TaskRepositoryFile")
	}

	tasksReturned, err := repository.GetAllTasks()
	if err != nil {
		t.Fatalf("expect GetAllTasks call to return no errors, got \"%s\"", err)
	}

	if len(tasksReturned) != len(tasks) {
		t.Fatalf("expect %d tasks to be returned, got %d", len(tasks), len(tasksReturned))
	}
}

func Test_DeleteTask(t *testing.T) {
	fileName := "Test_DeleteTask"
	defer removeTestFile(fileName)

	tasks := newTasks(2)
	addTasksToFileOrFail(tasks, fileName, t)

	repository, err := NewTaskRepositoryFile(fileName)
	if err != nil {
		t.Fatalf("failed to create TaskRepositoryFile: %s", err)
	}

	taskToDelete := tasks[1]
	if err := repository.DeleteTask(taskToDelete.Id); err != nil {
		t.Fatalf("expected DeleteTask call to return no errors, got %s", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("failed to open file %s: %s", fileName, err)
	}
	defer file.Close()

	var tasksInFile []model.Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasksInFile); err != nil {
		t.Errorf("failed to parse json from file %s: \"%v\"", fileName, err)
	}

	for _, task := range tasksInFile {
		if task.Id == taskToDelete.Id {
			t.Errorf("expected task %d to be deleted", taskToDelete.Id)
		}
	}

}

func newTasks(numberOfTasks int) []model.Task {
	tasks := make([]model.Task, numberOfTasks)
	for i := range numberOfTasks {
		task := model.Task{
			Id:          i + 1,
			Description: fmt.Sprintf("Task %d", i),
			Status:      0,
			CreatedAt:   model.DateTime(time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)),
			UpdatedAt:   model.DateTime(time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		}
		tasks[i] = task
	}
	return tasks
}

func addTasksToFileOrFail(tasks []model.Task, path string, t *testing.T) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file %s: %s", path, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if _, err = writer.WriteString("[\n"); err != nil {
		t.Fatalf("failed to write to buffer: %v", err)
	}

	for i, task := range tasks {
		jsonBytes, err := json.Marshal(&task)
		if err != nil {
			t.Fatalf("failed to serialize task: %v", err)
		}

		if _, err = writer.Write(jsonBytes); err != nil {
			t.Fatalf("failed to write to buffer: %v", err)
		}

		if i != len(tasks)-1 {
			if _, err = writer.WriteString(",\n"); err != nil {
				t.Fatalf("failed to write newline in buffer: %v", err)
			}
		}

	}

	if _, err = writer.WriteString("\n]"); err != nil {
		t.Fatalf("failed to write to buffer: %v", err)
	}

	if err = writer.Flush(); err != nil {
		t.Fatalf("failed to flush buffer: %v", err)
	}
}

func removeTestFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		panic(fmt.Errorf("failed to remove file %s: %w", fileName, err))
	}
}
