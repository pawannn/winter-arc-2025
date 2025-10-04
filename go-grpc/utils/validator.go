package utils

import (
	"net/mail"
	"regexp"
)

func ValidateName(name string) bool {
	re := regexp.MustCompile(`^[A-Za-z ]+$`)
	return re.MatchString(name)
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
