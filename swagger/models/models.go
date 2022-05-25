package models

type User struct {
	UserID int
}

type UserRepository interface {
	FindByLoginAndPwd(Login, Password string) (*User, error)
	FindLogin(Login string) error
	AddNewUser(Login, Password, FirstName, LastName, Email string, SocialMediaLinks []string) error
}
