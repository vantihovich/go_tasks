package validators

func ValidateRegistrationRequest(login, password, firstName, lastName, email string) bool {
	if login == "" || password == "" || firstName == "" || lastName == "" || email == "" {
		return false
	}
	return true
}

func ValidateChangePasswordRequest(newPassword, passwordConfirm string) bool {
	if newPassword == "" || newPassword != passwordConfirm {
		return false
	}
	return true
}

func ValidateWorldCoinIndexRequest(tickersList []string, fiat string) bool {
	if len(tickersList) == 0 || fiat == "" {
		return false
	}
	return true
}
