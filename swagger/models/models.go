package models

import "context"

type User struct {
	UserID    int
	UserLogin string
	Admin     bool
	Active    bool
}

type UserRepository interface {
	FindByLoginAndPwd(ctx context.Context, login, password string) (*User, error)
	CheckIfLoginExists(ctx context.Context, login string) (bool, error)
	AddNewUser(ctx context.Context, login, password, firstName, lastName, email string, socialMediaLinks []string) error
	GetAdminAttrUserLogin(ctx context.Context, userID int) (*User, error)
	DeactivateUser(ctx context.Context, userLogin string) error
	CheckIfUserActive(ctx context.Context, login string) (bool, error)
}
