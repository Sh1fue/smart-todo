package service

import (
	"context"
	"errors"
	"smart-todo/internal/domain"
	"smart-todo/internal/repository"
	"smart-todo/pkg/jwt"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, ErrUserExists
	}
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, ErrUserExists
	}

	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	token, err := s.jwtManager.Generate(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	{

		user, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return nil, ErrInvalidCredentials
		}

		if !user.CheckPassword(req.Password) {
			return nil, ErrInvalidCredentials
		}

		token, err := s.jwtManager.Generate(user.ID, user.Email)
		if err != nil {
			return nil, err
		}

		return &AuthResponse{
			Token: token,
			User:  user,
		}, nil
	}
}
