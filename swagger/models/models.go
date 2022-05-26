package models

import "context"

type User struct {
	UserID int
}

type UserRepository interface {
	FindByLoginAndPwd(ctx context.Context, login, password string) (*User, error)
	CheckIfLoginExists(ctx context.Context, login string) (bool, error)
	AddNewUser(ctx context.Context, login, password, firstName, lastName, email string, socialMediaLinks []string) error
}
