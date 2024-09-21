package repository

import (
	"fmt"
	"os"
	"testing"
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

	if repository.lineToWrite == 0 {
		t.Errorf("expect line to write to not be 0 but was %d", repository.lineToWrite)
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

	if repository.lineToWrite != 1 {
		t.Errorf("expect line to write to be 1 but was %d", repository.lineToWrite)
	}
}

func removeTestFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		panic(fmt.Errorf("failed to remove file %s: %w", fileName, err))
	}
}
