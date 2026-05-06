package handlers

import (
	"net/http"
	"todo_list/internal/db"
	"todo_list/internal/middleware"
	"todo_list/internal/models"
	"todo_list/web/templates"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	db *db.Database
}

func NewTaskHandler(database *db.Database) *TaskHandler {
	return &TaskHandler{db: database}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	if tenantID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tasks, err := h.db.GetTasks(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "Error al obtener tareas", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	w.Header().Set("Content-Type", "text/html")
	templates.TasksPage(tasks).Render(r.Context(), w)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Título requerido", http.StatusBadRequest)
		return
	}

	task, err := h.db.CreateTask(r.Context(), tenantID, title)
	if err != nil {
		http.Error(w, "Error al crear tarea", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.TaskItem(*task).Render(r.Context(), w)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	completed := r.FormValue("completed") == "on"

	task, err := h.db.UpdateTask(r.Context(), taskID, tenantID, completed)
	if err != nil {
		http.Error(w, "Error al actualizar tarea", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.TaskItem(*task).Render(r.Context(), w)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteTask(r.Context(), taskID, tenantID); err != nil {
		http.Error(w, "Error al eliminar tarea", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
