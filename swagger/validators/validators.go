package validators

func ValidateRegistrationRequest(login, password, firstName, lastName, email string) bool {
	if login == "" || password == "" || firstName == "" || lastName == "" || email == "" {
		return false
	}
	return true
}
