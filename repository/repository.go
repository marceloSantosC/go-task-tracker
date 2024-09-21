package repository

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type TaskRepositoryFile struct {
	path        string
	lineToWrite int
}

func NewTaskRepositoryFile(path string) (TaskRepositoryFile, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic(fmt.Errorf("failed to read file: %w", err))
	}
	defer file.Close()

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		newFile, err := os.Create(path)
		if err != nil {
			panic(fmt.Errorf("failed to create file %s: %w", path, err))
		}
		defer newFile.Close()

		if _, err = newFile.WriteString("[\n\n]"); err != nil {
			panic(fmt.Errorf("failed to write initial data to file %s: %w", path, err))
		}

		return TaskRepositoryFile{path: path, lineToWrite: 1}, nil
	}

	return TaskRepositoryFile{path: path}, nil

}
