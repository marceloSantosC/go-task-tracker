package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-task-tracker/model"
	"os"
	"regexp"
	"sync"
	"time"
)

type TaskRepositoryFile struct {
	path       string
	offset     int64
	sequenceId int
	mutex      sync.Mutex
}

const (
	firstLineValue = "[\n"
	lastLineValue  = "\n]"
	filePerm       = 0600
)

func NewTaskRepositoryFile(path string) (TaskRepositoryFile, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return TaskRepositoryFile{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return TaskRepositoryFile{}, fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.Size() == 0 {
		if _, err := file.WriteString(firstLineValue + lastLineValue); err != nil {
			return TaskRepositoryFile{}, fmt.Errorf("failed to initialize file: %w", err)
		}
		return TaskRepositoryFile{path: path, offset: int64(len(firstLineValue)), sequenceId: 0}, nil
	}

	lastLineBytes := int64(len(lastLineValue))
	offset := fileInfo.Size() - lastLineBytes

	sequenceId, err := loadSequenceId(file)
	if err != nil {
		panic(fmt.Errorf("failed to create sequence id %s: %w", path, err))
	}

	return TaskRepositoryFile{path: path, offset: offset, sequenceId: sequenceId}, nil
}

func loadSequenceId(file *os.File) (int, error) {
	decoder := json.NewDecoder(file)
	var tasks []model.Task
	if err := decoder.Decode(&tasks); err != nil {
		return 0, fmt.Errorf("failed to decode tasks in file %s: %w", file.Name(), err)
	}

	if len(tasks) == 0 {
		return 0, nil
	}

	return tasks[len(tasks)-1].Id, nil
}

func (r *TaskRepositoryFile) AddTask(task model.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	file, err := os.OpenFile(r.path, os.O_RDWR, filePerm)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", r.path, err)
	}
	defer file.Close()

	r.sequenceId++
	task.Id = r.sequenceId

	b, err := json.Marshal(&task)
	if err != nil {
		return fmt.Errorf("failed to marshal task %d: %w", task.Id, err)
	}

	stringJson := string(b)
	if r.offset != int64(len(firstLineValue)) {
		stringJson = fmt.Sprintf(",\n%s%s", stringJson, lastLineValue)
	} else {
		stringJson = fmt.Sprintf("%s%s", stringJson, lastLineValue)
	}

	writtenBytes, err := file.WriteAt([]byte(stringJson), r.offset)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	r.offset += int64(writtenBytes)
	return nil
}

func (r *TaskRepositoryFile) UpdateTask(id int, updatedTask model.UpdateTask) error {

	file, err := os.OpenFile(r.path, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", r.path, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	regexExpr := fmt.Sprintf("(\"Id\":%d)", id)

	for scanner.Scan() {
		line := scanner.Text()
		if matches, _ := regexp.MatchString(regexExpr, line); matches {
			var task model.Task
			if err = json.Unmarshal([]byte(line[:len(line)-1]), &task); err != nil {
				return fmt.Errorf("failed to unmarshal json: %w", err)
			}
			task.Description = updatedTask.Description
			task.Status = updatedTask.Status
			task.UpdatedAt = model.DateTime(time.Now())

			jsonBytes, err := json.Marshal(&task)
			if err != nil {
				return fmt.Errorf("failed to marshal json: %w", err)
			}
			line = string(jsonBytes) + line[len(line)-1:]
		}
		lines = append(lines, line)
	}

	if err := truncateAndWrite(file, lines); err != nil {
		return fmt.Errorf("failed to update task %d: %w", id, err)
	}

	return nil
}

func (r *TaskRepositoryFile) GetAllTasks() ([]model.Task, error) {
	file, err := os.Open(r.path)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var tasks []model.Task
	if err = decoder.Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil

}

func (r *TaskRepositoryFile) DeleteTask(id int) error {
	file, err := os.OpenFile(r.path, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to retrieve tasks: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	regexExpr := fmt.Sprintf("(\"Id\":%d)", id)

	taskToDeleteIndex := -1
	for scanner.Scan() {
		line := scanner.Text()
		if matches, _ := regexp.MatchString(regexExpr, line); matches {
			taskToDeleteIndex = len(lines)
		}
		lines = append(lines, line)
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	if taskToDeleteIndex == -1 {
		return fmt.Errorf("task with id %d does not exists", id)
	}

	if len(lines) > 1 && taskToDeleteIndex == len(lines)-2 {
		lineEndsWithComma := regexp.MustCompile(",$")
		lines[taskToDeleteIndex-1] = lineEndsWithComma.ReplaceAllString(lines[taskToDeleteIndex-1], "")
	}

	lines = append(lines[:taskToDeleteIndex], lines[taskToDeleteIndex+1:]...)

	if err := truncateAndWrite(file, lines); err != nil {
		return fmt.Errorf("failed to delete task %d: %w", id, err)
	}

	return nil

}

func truncateAndWrite(file *os.File, lines []string) error {
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek to beginning of file: %w", err)
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write in buffer: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}
