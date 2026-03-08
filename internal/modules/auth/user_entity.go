package auth

type User struct {
	ID           string
	Username     string
	PasswordHash string
	FullName     string
}
