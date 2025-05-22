package domain

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"` // не возвращаем хеш пароля наружу
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Phone     string    `json:"phone"`
	IsAdmin   bool      `json:"is_admin"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=3"`
	Surname  string `json:"surname" validate:"required,min=3"`
	Phone    string `json:"phone" validate:"required,min=6"`
}

type UserAuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=6"`
	Name     string `json:"name,omitempty" validate:"omitempty,min=3"`
	Surname  string `json:"surname,omitempty" validate:"omitempty,min=3"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=6"`
}
