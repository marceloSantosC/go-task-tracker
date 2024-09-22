package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-task-tracker/model"
	"io"
	"io/fs"
	"os"
)

type TaskRepositoryFile struct {
	path   string
	offset int64
}

const firstLineValue = "[\n"
const lastLineValue = "\n]"

func NewTaskRepositoryFile(path string) (TaskRepositoryFile, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0700)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic(fmt.Errorf("failed to read file: %w", err))
	}
	defer file.Close()

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		file, err = os.Create(path)
		if err != nil {
			panic(fmt.Errorf("failed to create file %s: %w", path, err))
		}
		defer file.Close()

		if _, err = file.WriteString(firstLineValue + lastLineValue); err != nil {
			panic(fmt.Errorf("failed to write initial data to file %s: %w", path, err))
		}
	}

	totalBytes, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return TaskRepositoryFile{}, fmt.Errorf("failed to retrieve status for file %s: %w", path, err)
	}
	lastLineBytes := int64(len(lastLineValue))
	offset := totalBytes - lastLineBytes
	return TaskRepositoryFile{path: path, offset: offset}, nil
}

func (r TaskRepositoryFile) AddTask(task model.Task) error {
	file, err := os.OpenFile(r.path, os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("failed to add task, cound'nt open file %s: %w", r.path, err)
	}
	defer file.Close()

	b, err := json.Marshal(&task)
	if err != nil {
		return fmt.Errorf("failed to add task, coundn't marshal json: %w", err)
	}

	stringJson := string(b)
	if r.offset != int64(len(firstLineValue)) {
		stringJson = fmt.Sprintf(",\n%s%s", stringJson, lastLineValue)
	} else {
		stringJson = fmt.Sprintf("%s%s", stringJson, lastLineValue)
	}

	writtenBytes, err := file.WriteAt([]byte(stringJson), r.offset)
	if err != nil {
		return fmt.Errorf("failed to add task, coundn't write in file %s: %w", r.path, err)
	}
	r.offset += int64(writtenBytes)
	return nil
}
