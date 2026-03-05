package domain

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrUsernameExists  = errors.New("username already exists")
	ErrTaskNotFound    = errors.New("task not found")
	ErrInvalidPriority = errors.New("invalid priority")
	ErrInvalidStatus   = errors.New("invalid status")
)