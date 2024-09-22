package main

import (
	"go-task-tracker/repository"
)

func main() {

	filename := "task_list"
	tr, _ := repository.NewTaskRepositoryFile(filename)
	_ = tr
}
