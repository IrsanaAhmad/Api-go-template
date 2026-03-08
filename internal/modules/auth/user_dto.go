package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,max=100"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	User UserDTO `json:"user"`
}

type RegisterResponse struct {
	User UserDTO `json:"user"`
}

type RefreshResponse struct {
	User UserDTO `json:"user"`
}
