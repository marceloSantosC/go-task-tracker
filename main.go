package main

import (
	"fmt"
	"go-task-tracker/repository"
	"go-task-tracker/service"
	"log"
	"log/slog"
	"os"
	"time"
)

func main() {

	filename := "task_list.json"
	repo, err := repository.NewTaskRepositoryFile(filename)
	if err != nil {
		log.Panicf("failed to start app: %s", err)
	}

	file, err := os.Create(fmt.Sprintf("logs_%d", time.Now().UnixMilli()))
	if err != nil {
		log.Panicf("failed to start app: %s", err)
	}

	l := slog.New(slog.NewTextHandler(file, nil))

	l.Info("Initialized app using file.", slog.String("file", filename))
	service.NewTaskService(&repo, l)

	fmt.Println("# GO Task Tracker #")
}
