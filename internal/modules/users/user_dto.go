package users

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserDTO `json:"user"`
}
