package common

import (
	"net/mail"
	"regexp"
	"strings"
)

var (
	domainRegex = regexp.MustCompile(
		`^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$`,
	)
	handleRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,31}$`)
)

func ValidateEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return addr.Address == email
}

func ValidatePassword(password string) bool {
	if len(password) < 12 || len(password) > 64 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func ValidateDomain(domain string) bool {
	return domainRegex.MatchString(domain)
}

func ValidateHandle(handle string) bool {
	return handleRegex.MatchString(handle)
}

func ValidateCountryCode(code string) bool {
	if len(code) != 3 {
		return false
	}
	for _, r := range code {
		if !('A' <= r && r <= 'Z') {
			return false
		}
	}
	return true
}

func ValidateCurrency(currency string) bool {
	if len(currency) != 3 {
		return false
	}
	for _, r := range currency {
		if !('A' <= r && r <= 'Z') {
			return false
		}
	}
	return true
}
