package repository

import (
	"encoding/json"
	"fmt"
	"go-task-tracker/model"
	"os"
	"sync"
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
