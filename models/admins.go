package models

import (
	"time"

	"github.com/google/uuid"
)

type Admins struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" `
}

type CreateAdmin struct {
	Name string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateAdmin struct {
	Name string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type GetListAdmins struct {
	Admins []Admins `json:"admins"`
	Count  int     `json:"count"`
}
