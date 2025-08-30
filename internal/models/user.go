package models

// User struct keeps information about user
type User struct {
	ID       string
	Login    string
	Password []byte
	Salt     []byte
}
