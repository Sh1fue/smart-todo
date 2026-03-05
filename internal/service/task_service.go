package service

import (
	"context"
	"smart-todo/internal/domain"
	"smart-todo/internal/repository"
	"time"
)

type TaskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

type CreateTaskRequest struct {
	UserID      int             `json:"user_id"`
	Title       string          `json:"title" validate:"required,min=1,max=200"`
	Description string          `json:"description"`
	Priority    domain.Priority `json:"priority"`
	DueDate     time.Time       `json:"due_date"`
	Tags        []string        `json:"tags"`
}

type UpdateTaskRequest struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Priority    domain.Priority `json:"priority"`
	Status      domain.Status   `json:"status"`
	DueDate     time.Time       `json:"due_date"`
	Tags        []string        `json:"tags"`
}

func (s *TaskService) Create(ctx context.Context, req *CreateTaskRequest) (*domain.Task, error) {
	task := &domain.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
		Status:      domain.StatusActive,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *TaskService) GetUserTasks(ctx context.Context, userID int) ([]*domain.Task, error) {
	return s.taskRepo.GetByUserID(ctx, userID)
}

func (s *TaskService) Update(ctx context.Context, req *UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Обновляем только переданные поля
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Priority != 0 {
		task.Priority = req.Priority
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if !req.DueDate.IsZero() {
		task.DueDate = req.DueDate
	}
	if req.Tags != nil {
		task.Tags = req.Tags
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id int) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *TaskService) MarkAsDone(ctx context.Context, id int) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	task.Status = domain.StatusDone
	return s.taskRepo.Update(ctx, task)
}

// Умные методы
func (s *TaskService) GetOverdueTasks(ctx context.Context, userID int) ([]*domain.Task, error) {
	return s.taskRepo.GetOverdue(ctx, userID)
}

func (s *TaskService) GetTasksByPriority(ctx context.Context, userID int, priority domain.Priority) ([]*domain.Task, error) {
	return s.taskRepo.GetByPriority(ctx, userID, priority)
}

func (s *TaskService) GetStats(ctx context.Context, userID int) (*domain.TaskStats, error) {
	return s.taskRepo.GetStats(ctx, userID)
}
