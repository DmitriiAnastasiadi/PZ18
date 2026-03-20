package http

import (
	"encoding/json"

	"net/http"

	"tech-ip-sem2/services/tasks/internal/service"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Title string `json:"title"`

		Description string `json:"description,omitempty"`

		DueDate string `json:"due_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(w, "Bad request", 400)

		return

	}

	task := service.CreateTask(req.Title, req.Description, req.DueDate)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(201)

	json.NewEncoder(w).Encode(task)

}

func ListTasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks := service.GetTasks()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tasks)

}

func GetTaskHandler(w http.ResponseWriter, r *http.Request, id string) {

	task, ok := service.GetTask(id)

	if !ok {

		http.NotFound(w, r)

		return

	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(task)

}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, id string) {

	var updates map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {

		http.Error(w, "Bad request", 400)

		return

	}

	task, ok := service.UpdateTask(id, updates)

	if !ok {

		http.NotFound(w, r)

		return

	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(task)

}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, id string) {

	if !service.DeleteTask(id) {

		http.NotFound(w, r)

		return

	}

	w.WriteHeader(204)

}
