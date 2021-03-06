package models

import (
	"context"
	"time"
)

type User struct {
	ID       int
	Login    string
	Password string
	RoleID   string
	Active   bool
}

type UserRepository interface {
	FindByLogin(ctx context.Context, userLogin string) (*User, error)
	FindByID(ctx context.Context, userID int) (*User, error)
	CheckIfLoginExists(ctx context.Context, login string) (bool, error)
	AddNewUser(ctx context.Context, login, password, firstName, lastName, email string, socialMediaLinks []string) error
	GetAdminAttrUserLogin(ctx context.Context, userID int) (*User, error)
	DeactivateUser(ctx context.Context, userLogin string) (bool, error)
	ChangePassword(ctx context.Context, userID int, newPassword string) error
	WriteSecretToDBForPasswordRecovery(ctx context.Context, email, secret string) (bool, error)
	CheckSecretChangePassword(ctx context.Context, secret string, newPassword string) error
}

var roleAdmin = "1"

func (u *User) IsAdmin() bool {
	return u.RoleID == roleAdmin
}

type WorldCoinIndexHandlerResponse struct {
	Label      string
	Name       string
	Price      float32
	Volume_24h float32
	Timestamp  time.Time
}

type WorldCoinIndexHandlerResponseArray struct {
	Response []WorldCoinIndexHandlerResponse
}

type WCISource interface {
	LoadTickers(requestList []string, fiat string) (response WorldCoinIndexHandlerResponseArray, err error)
}
