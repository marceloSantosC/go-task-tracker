package repository

import (
	"errors"
	"fmt"
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
	file, err := os.OpenFile(path, os.O_RDONLY, 0777)
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
