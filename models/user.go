package models

import "fmt"

//User implements a generic user model
type User struct {
	ID       int64
	UserName string `sql:",unique"`
	Email    string
	Password string
	IsAdmin  bool
}

func (u User) String() string {
	text := fmt.Sprintf("%s %s %s", u.UserName, u.Email, u.Password)
	return text
}
