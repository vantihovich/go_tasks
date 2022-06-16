package validators

func ValidateRegistrationRequest(login, password, firstName, lastName, email string) bool {
	if login == "" || password == "" || firstName == "" || lastName == "" || email == "" {
		return false
	}
	return true
}

func ValidateChangingPWD(newPWD, newPWDConfirm string) bool {
	if newPWD != newPWDConfirm || newPWD == "" || newPWDConfirm == "" {
		return false
	}
	return true
}
