package models

import "github.com/google/uuid"

type Categories struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CreateCategories struct {
	Name string `json:"name"`
}

type UpdateCategories struct {
	Name string `json:"name"`
}

type GetListCategories struct {
	Categories []*Categories `json:"categories"`
	Count      int           `josn:"count"`
}
