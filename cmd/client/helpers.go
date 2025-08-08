package main

import (
	"fmt"
	"unicode/utf8"
)

func validateLogin(login string) error {
	if utf8.RuneCountInString(login) < 3 {
		return fmt.Errorf("login must be at least 3 characters")
	}
	if utf8.RuneCountInString(login) > 255 {
		return fmt.Errorf("login must be at most 255 characters")
	}

	return nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if utf8.RuneCountInString(password) > 36 {
		return fmt.Errorf("password must be at most 36 characters")
	}
	if len(password) > 72 {
		return fmt.Errorf("password is %d characters but is longer than 72 bytes", utf8.RuneCountInString(password))
	}

	return nil
}
