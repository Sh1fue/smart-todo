package dto

import (
	"smart-todo/internal/domain"
	"time"
)

type CreateTaskRequest struct {
	Title       string          `json:"title" validate:"required,min=1,max=200"`
	Description string          `json:"description"`
	Priority    domain.Priority `json:"priority" validate:"min=0,max=3"`
	DueDate     time.Time       `json:"due_date"`
	Tags        []string        `json:"tags"`
}

type UpdateTaskRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Priority    domain.Priority `json:"priority"`
	Status      domain.Status   `json:"status"`
	DueDate     time.Time       `json:"due_date"`
	Tags        []string        `json:"tags"`
}

type TaskResponse struct {
	ID          int             `json:"id"`
	UserID      int             `json:"user_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Priority    string          `json:"priority"`
	PriorityID  domain.Priority `json:"priority_id"`
	Status      string          `json:"status"`
	StatusID    domain.Status   `json:"status_id"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	DueDate     string          `json:"due_date"`
	CompletedAt string          `json:"completed_at"`
	Tags        []string        `json:"tags"`
	IsOverdue   bool            `json:"is_overdue"`
}

type TaskStatsResponse struct {
	Total      int               `json:"total"`
	Completed  int               `json:"completed"`
	Overdue    int               `json:"overdue"`
	ByPriority map[string]int    `json:"by_priority"`
	ByStatus   map[string]int    `json:"by_status"`
}

func ToTaskResponse(task *domain.Task) *TaskResponse {
	if task == nil {
		return nil
	}
	
	return &TaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority.String(),
		PriorityID:  task.Priority,
		Status:      string(task.Status),
		StatusID:    task.Status,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
		DueDate:     task.DueDate.Format(time.RFC3339),
		CompletedAt: task.CompletedAt.Format(time.RFC3339),
		Tags:        task.Tags,
		IsOverdue:   task.IsOverdue(),
	}
}

func ToTaskResponses(tasks []*domain.Task) []*TaskResponse {
	responses := make([]*TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = ToTaskResponse(task)
	}
	return responses
}

func ToTaskStatsResponse(stats *domain.TaskStats) *TaskStatsResponse {
	if stats == nil {
		return nil
	}
	
	byPriority := make(map[string]int)
	for p, count := range stats.ByPriority {
		byPriority[p.String()] = count
	}
	
	byStatus := make(map[string]int)
	for s, count := range stats.ByStatus {
		byStatus[string(s)] = count
	}
	
	return &TaskStatsResponse{
		Total:      stats.Total,
		Completed:  stats.Completed,
		Overdue:    stats.Overdue,
		ByPriority: byPriority,
		ByStatus:   byStatus,
	}
}