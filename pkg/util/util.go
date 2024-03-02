package util

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2}`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) (string, bool) {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
		// specialRunes = "!@#$%^&*"
	)

	if len(password) < 8 {
		return "Password must contain at least 8 characters", false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "Password must contain an upper case character", false
	}
	if !hasLower {
		return "Password must contain a lower case character", false
	}
	if !hasNumber {
		return "Password must contain a number", false
	}

	if !hasSpecial {
		return "Password must contain a special character", false
	}

	return "", true

}
