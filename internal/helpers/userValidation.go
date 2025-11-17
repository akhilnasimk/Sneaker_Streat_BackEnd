package helpers

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateUserInput(username, email, password, phone string) error {

	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	phone = strings.TrimSpace(phone)

	// Required fields
	if username == "" || email == "" || password == "" || phone == "" { //checking if empty
		return errors.New("username, email, password, and phone are required")
	}

	// Username validation
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters long")
	}

	// Email format validation
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	if ok, _ := regexp.MatchString(emailRegex, email); !ok {
		return errors.New("invalid email format")
	}

	// Password validation
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Strong password check (optional but recommended)
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) ||
		!regexp.MustCompile(`[a-z]`).MatchString(password) ||
		!regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain uppercase, lowercase, and a number")
	}
	if len(phone) != 10 {
		return errors.New("phone number must be 10 number ")
	}

	return nil
}
