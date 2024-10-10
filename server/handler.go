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
	http.HandleFunc("PUT /tasks/{id}", h.HandleUpdateTask)
	http.HandleFunc("DELETE /tasks/{id}", h.HandleDeleteTask)
	return h
}

func (h TaskHandler) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var task model.CreateTask
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&task); err != nil {
		h.log.Error("failed to process request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.AddTask(task); err != nil {
		h.log.Error("failed to process request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	response, err := h.service.GetTasks(model.TaskStatus(status), description)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to get tasks: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonRes, err := json.Marshal(&response)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to marshal json: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonRes); err != nil {
		h.log.Error(fmt.Sprintf("error when writing http response: %s", err))
	}
}

func (h TaskHandler) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(fmt.Sprintf("invalid path variable id with value %s", r.PathValue("id")))
		w.WriteHeader(http.StatusBadRequest)
	}

	var task model.UpdateTask
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&task); err != nil {
		h.log.Error("failed to process request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTask(id, task); err != nil {
		h.log.Error("failed to process request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h TaskHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(fmt.Sprintf("invalid path variable id with value %s", r.PathValue("id")))
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := h.service.DeleteTask(id); err != nil {
		h.log.Error("failed to process request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
