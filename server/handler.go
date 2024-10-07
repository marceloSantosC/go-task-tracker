package server

import (
	"encoding/json"
	"fmt"
	"go-task-tracker/model"
	"go-task-tracker/service"
	"log/slog"
	"net/http"
	"strconv"
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
	http.HandleFunc("GET /tasks", h.HandleGetTasks)
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

func (h TaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	statusFilter := r.URL.Query().Get("status")
	description := r.URL.Query().Get("description")

	status := -1
	if statusFilter != "" {
		var err error
		if status, err = strconv.Atoi(statusFilter); err != nil {
			h.log.Info(fmt.Sprintf("input %s is invalid for query param status", statusFilter))
			w.WriteHeader(400)
		}
	}

	response, err := h.service.GetTasks(model.TaskStatus(status), description)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to get tasks: %s", err))
		w.WriteHeader(500)
	}

	jsonRes, err := json.Marshal(&response)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to marshal json: %s", err))
		w.WriteHeader(500)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonRes); err != nil {
		h.log.Error(fmt.Sprintf("error when writing http response: %s", err))
	}
}
