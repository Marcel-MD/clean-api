package models

import "gorm.io/datatypes"

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	Base

	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`

	Roles datatypes.JSONSlice[string] `json:"roles"`
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}

type RegisterUser struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type Token struct {
	Token string `json:"token"`
}
