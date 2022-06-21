package validators

func ValidateRegistrationRequest(login, password, firstName, lastName, email string) bool {
	if login == "" || password == "" || firstName == "" || lastName == "" || email == "" {
		return false
	}
	return true
}

func ValidateChangingPWD(newPassword, passwordConfirm string) bool {
	if newPassword == "" || newPassword != passwordConfirm {
		return false
	}
	return true
}
