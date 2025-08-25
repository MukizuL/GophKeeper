package models

type User struct {
	ID       string
	Login    string
	Password []byte
	Salt     []byte
}
