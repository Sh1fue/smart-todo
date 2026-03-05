package domain

import (
	"time"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "Низкий"
	case PriorityMedium:
		return "Средний"
	case PriorityHigh:
		return "Высокий"
	case PriorityCritical:
		return "Критичный"
	default:
		return "Неизвестно"
	}
}

type Status string

const (
	StatusActive  Status = "Активна"
	StatusDone    Status = "Завершена" // Исправил: "Завершенная" -> "Завершена"
	StatusOverdue Status = "Просрочена" // Исправил: "Просроченна" -> "Просрочена"
)

type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`      // 👈 ДОБАВИТЬ: связь с пользователем
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`   // 👈 ИСПРАВИТЬ: CreateAt -> CreatedAt
	UpdatedAt   time.Time `json:"updated_at"`   // 👈 ДОБАВИТЬ: время обновления
	DueDate     time.Time `json:"due_date"`
	CompletedAt time.Time `json:"completed_at"` // 👈 ИСПРАВИТЬ: DoneAt -> CompletedAt (для консистентности)
	Tags        []string  `json:"tags"`
}
	type TaskStats struct {
	Total      int                `json:"total"`
	Completed  int                `json:"completed"`
	Overdue    int                `json:"overdue"`
	ByPriority map[Priority]int   `json:"by_priority"`
	ByStatus   map[Status]int     `json:"by_status"`
}

func (t *Task) IsOverdue() bool {
	if t.Status == StatusDone {
		return false
	}
	return !t.DueDate.IsZero() && t.DueDate.Before(time.Now())
}

func (t *Task) UpdateStatus() {
	if t.Status == StatusDone {
		return
	}
	if t.IsOverdue() {
		t.Status = StatusOverdue
	} else {
		t.Status = StatusActive
	}
}