package models

type Email string

type User struct {
	Email Email
}

func NewUser(email Email) *User {
	return &User{Email: email}
}

func (u *User) GetEmail() Email {
	return u.Email
}
