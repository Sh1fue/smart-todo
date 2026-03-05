package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/internal/api/dto"
	"smart-todo/internal/api/middleware"
	"smart-todo/internal/domain"
	"smart-todo/internal/service"
	"strconv"
	"strings"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask создает новую задачу
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из контекста (установлен middleware)
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: добавить валидацию

	task, err := h.taskService.Create(r.Context(), &service.CreateTaskRequest{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
	})

	if err != nil {
		http.Error(w, "failed to create task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToTaskResponse(task))
}

// GetTasks возвращает все задачи пользователя
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсим query параметры для фильтрации
	query := r.URL.Query()
	priority := query.Get("priority")
	status := query.Get("status")
	overdue := query.Get("overdue")

	tasks, err := h.taskService.GetUserTasks(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Фильтруем задачи по параметрам
	var filteredTasks []*domain.Task
	for _, task := range tasks {
		// Фильтр по приоритету
		if priority != "" {
			p, _ := strconv.Atoi(priority)
			if int(task.Priority) != p {
				continue
			}
		}
		// Фильтр по статусу
		if status != "" && string(task.Status) != status {
			continue
		}
		// Фильтр по просроченности
		if overdue == "true" && !task.IsOverdue() {
			continue
		}
		if overdue == "false" && task.IsOverdue() {
			continue
		}
		filteredTasks = append(filteredTasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToTaskResponses(filteredTasks))
}

// GetTask возвращает задачу по ID
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetByID(r.Context(), taskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// Проверяем, что задача принадлежит пользователю
	if task.UserID != userID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToTaskResponse(task))
}

// UpdateTask обновляет задачу
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	// Проверяем существование задачи и права доступа
	existing, err := h.taskService.GetByID(r.Context(), taskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	if existing.UserID != userID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.Update(r.Context(), &service.UpdateTaskRequest{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
	})

	if err != nil {
		http.Error(w, "failed to update task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToTaskResponse(task))
}

// DeleteTask удаляет задачу
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	// Проверяем права доступа
	existing, err := h.taskService.GetByID(r.Context(), taskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	if existing.UserID != userID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if err := h.taskService.Delete(r.Context(), taskID); err != nil {
		http.Error(w, "failed to delete task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MarkAsDone отмечает задачу как выполненную
func (h *TaskHandler) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	// Проверяем права доступа
	existing, err := h.taskService.GetByID(r.Context(), taskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	if existing.UserID != userID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if err := h.taskService.MarkAsDone(r.Context(), taskID); err != nil {
		http.Error(w, "failed to mark task as done: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "task marked as done"})
}

// GetStats возвращает статистику по задачам
func (h *TaskHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	stats, err := h.taskService.GetStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToTaskStatsResponse(stats))
}

// GetOverdue возвращает просроченные задачи
func (h *TaskHandler) GetOverdue(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.taskService.GetOverdueTasks(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get overdue tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToTaskResponses(tasks))
}
