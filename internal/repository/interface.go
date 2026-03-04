package repository

import (
	"context"
	"smart-todo/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id int) (*domain.Task, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id int) error
	// Умные методы
	GetOverdue(ctx context.Context, userID int) ([]*domain.Task, error)
	GetByPriority(ctx context.Context, userID int, priority domain.Priority) ([]*domain.Task, error)
	GetStats(ctx context.Context, userID int) (*TaskStats, error)
}

type TaskStats struct {
	Total     int            `json:"total"`
	Completed int            `json:"completed"`
	Overdue   int            `json:"overdue"`
	ByPriority map[domain.Priority]int `json:"by_priority"`
}