package entity

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-" validate:"required,min=8"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(name, email, password string) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) Update(name, email string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if email == "" {
		return errors.New("email is required")
	}

	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}
