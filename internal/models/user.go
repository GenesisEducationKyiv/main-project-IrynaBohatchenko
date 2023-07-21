package models

type Email string

type User struct {
	email Email
}

func NewUser(email Email) *User {
	return &User{email: email}
}

func (u *User) GetEmail() Email {
	return u.email
}
