package models

type User struct {
	UserID string
}

type UserRepository interface {
	FindByLoginAndPwd(Login, Password string) (*User, error)
}
