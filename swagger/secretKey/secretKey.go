package secretKey

import (
	str "github.com/xyproto/randomstring"
)

func NewRandomString(length int) (secret string) {
	return str.HumanFriendlyEnglishString(length)
}
