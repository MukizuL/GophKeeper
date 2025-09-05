package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) < 16+3 {
		return fmt.Errorf("CCN is too short")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash, and it should be in the 2nd index (3rd character)
	if len(s) != 5 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	if len(s) != 3 {
		return fmt.Errorf("CVV is invalid")
	}
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

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
