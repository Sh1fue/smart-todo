package memory

import (
	"context"
	"sync"
	"time"
	"smart-todo/internal/domain"
)

type TaskRepository struct {
	mu        sync.RWMutex
	tasks     map[int]*domain.Task
	userTasks map[int][]int // user_id -> []task_id
	lastID    int
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks:     make(map[int]*domain.Task),
		userTasks: make(map[int][]int),
		lastID:    0,
	}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastID++
	task.ID = r.lastID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = domain.StatusActive

	r.tasks[task.ID] = task
	r.userTasks[task.UserID] = append(r.userTasks[task.UserID], task.ID)

	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, domain.ErrTaskNotFound
	}

	// Обновляем статус перед возвратом
	task.UpdateStatus()
	return task, nil
}

func (r *TaskRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taskIDs, exists := r.userTasks[userID]
	if !exists {
		return []*domain.Task{}, nil
	}

	tasks := make([]*domain.Task, 0, len(taskIDs))
	for _, id := range taskIDs {
		if task, exists := r.tasks[id]; exists {
			// Копируем задачу, чтобы не изменять оригинал при обновлении статуса
			taskCopy := *task
			taskCopy.UpdateStatus()
			tasks = append(tasks, &taskCopy)
		}
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.tasks[task.ID]
	if !exists {
		return domain.ErrTaskNotFound
	}

	// Обновляем время завершения если статус изменился на Done
	if task.Status == domain.StatusDone && existing.Status != domain.StatusDone {
		task.CompletedAt = time.Now()
	}

	task.UpdatedAt = time.Now()
	task.CreatedAt = existing.CreatedAt // сохраняем оригинальное время создания
	r.tasks[task.ID] = task

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return domain.ErrTaskNotFound
	}

	// Удаляем из userTasks
	userTaskIDs := r.userTasks[task.UserID]
	for i, taskID := range userTaskIDs {
		if taskID == id {
			r.userTasks[task.UserID] = append(userTaskIDs[:i], userTaskIDs[i+1:]...)
			break
		}
	}

	delete(r.tasks, id)
	return nil
}

// Умные методы
func (r *TaskRepository) GetOverdue(ctx context.Context, userID int) ([]*domain.Task, error) {
	tasks, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	overdue := make([]*domain.Task, 0)
	for _, task := range tasks {
		if task.IsOverdue() {
			overdue = append(overdue, task)
		}
	}

	return overdue, nil
}

func (r *TaskRepository) GetByPriority(ctx context.Context, userID int, priority domain.Priority) ([]*domain.Task, error) {
	tasks, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	filtered := make([]*domain.Task, 0)
	for _, task := range tasks {
		if task.Priority == priority {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

func (r *TaskRepository) GetStats(ctx context.Context, userID int) (*domain.TaskStats, error) {
	tasks, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := &domain.TaskStats{
		Total:      len(tasks),
		Completed:  0,
		Overdue:    0,
		ByPriority: make(map[domain.Priority]int),
		ByStatus:   make(map[domain.Status]int),
	}

	for _, task := range tasks {
		if task.Status == domain.StatusDone {
			stats.Completed++
		}
		if task.IsOverdue() {
			stats.Overdue++
		}
		stats.ByPriority[task.Priority]++
		stats.ByStatus[task.Status]++
	}

	return stats, nil
}