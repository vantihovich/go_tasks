package models

import "context"

type User struct {
	ID     int
	Login  string
	RoleID string
	Active bool
}

type UserRepository interface {
	FindByLoginAndPwd(ctx context.Context, login, password string) (*User, error)
	CheckIfLoginExists(ctx context.Context, login string) (bool, error)
	AddNewUser(ctx context.Context, login, password, firstName, lastName, email string, socialMediaLinks []string) error
	GetAdminAttrUserLogin(ctx context.Context, userID int) (*User, error)
	DeactivateUser(ctx context.Context, userLogin string) (bool, error)
	ChangePWD(ctx context.Context, userLogin, newPassword string) error
}

var roleAdmin = "1"

func (u *User) IsAdmin() bool {
	return u.RoleID == roleAdmin
}
