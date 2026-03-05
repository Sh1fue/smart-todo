package memory

import (
	"context"
	"sync"
	"smart-todo/internal/domain"
)

type UserRepository struct {
	mu     sync.RWMutex
	users  map[int]*domain.User
	emailIndex map[string]int
	usernameIndex map[string]int
	lastID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:        make(map[int]*domain.User),
		emailIndex:   make(map[string]int),
		usernameIndex: make(map[string]int),
		lastID:       0,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем уникальность email
	if _, exists := r.emailIndex[user.Email]; exists {
		return domain.ErrEmailExists
	}

	// Проверяем уникальность username
	if _, exists := r.usernameIndex[user.Username]; exists {
		return domain.ErrUsernameExists
	}

	// Генерируем ID
	r.lastID++
	user.ID = r.lastID

	// Сохраняем пользователя
	r.users[user.ID] = user
	r.emailIndex[user.Email] = user.ID
	r.usernameIndex[user.Username] = user.ID

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.emailIndex[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return r.users[id], nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.usernameIndex[username]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return r.users[id], nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	// Обновляем индексы если изменился email
	if oldUser := r.users[user.ID]; oldUser.Email != user.Email {
		delete(r.emailIndex, oldUser.Email)
		r.emailIndex[user.Email] = user.ID
	}

	// Обновляем индексы если изменился username
	if oldUser := r.users[user.ID]; oldUser.Username != user.Username {
		delete(r.usernameIndex, oldUser.Username)
		r.usernameIndex[user.Username] = user.ID
	}

	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}

	delete(r.emailIndex, user.Email)
	delete(r.usernameIndex, user.Username)
	delete(r.users, id)

	return nil
}