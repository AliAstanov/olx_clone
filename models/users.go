package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId    uuid.UUID    `json:"user_id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `Json:"deleted_at"`
}

type CreateUserReq struct {
	Name     string `json:"name" validate:"required"`           // Majburiy maydon
	Email    string `json:"email" validate:"required,email"`    // Majburiy va email formatida
	Password string `json:"password" validate:"required,min=8"` // Majburiy va kamida 8 belgidan iborat
}

type UpdateUserReq struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type GetListUsers struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}
