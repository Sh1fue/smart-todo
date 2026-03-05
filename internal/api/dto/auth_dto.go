package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
