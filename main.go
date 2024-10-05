package main

import (
	"go-task-tracker/repository"
	"go-task-tracker/server"
	"go-task-tracker/service"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	filename := "task_list.json"
	repo, err := repository.NewTaskRepositoryFile(filename)
	if err != nil {
		log.Error("failed to start app", slog.String("error", err.Error()))
		panic(err)
	}

	log.Info("Initialized app using file.", slog.String("file", filename))
	s := service.NewTaskService(&repo, log)
	_ = server.NewTaskHandler(s, log)

	log.Info("Server started on port 8080")
	if err = http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
		panic(err)
	}

}
