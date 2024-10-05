package server

import (
	"encoding/json"
	"go-task-tracker/model"
	"go-task-tracker/service"
	"log/slog"
	"net/http"
)

type TaskHandler struct {
	service service.TaskService
	log     slog.Logger
}

func NewTaskHandler(service service.TaskService, log *slog.Logger) TaskHandler {
	h := TaskHandler{
		service: service,
		log:     *log,
	}
	http.HandleFunc("POST /tasks", h.HandlePostTask)
	return h
}

func (h TaskHandler) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var task model.CreateOrUpdateTask
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&task); err != nil {
		h.log.Error("failed to process request: ", err)
		w.WriteHeader(400)
		return
	}

	if err := h.service.AddTask(task); err != nil {
		h.log.Error("failed to process request: ", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
}
