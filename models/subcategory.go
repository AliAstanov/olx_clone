package models

import "github.com/google/uuid"

type Subcategories struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CategoryID uuid.UUID `json:"category_id"`
}

type CreateSubcategories struct {
	Name       string    `json:"name"`
	CategoryID uuid.UUID `json:"category_id"`
}

type UpdateSubcategories struct {
	Name       string    `json:"name"`
	CategoryID uuid.UUID `json:"category_id"`
}

type GetListSubcategories struct {
	Subcategories []Subcategories `json:"subcategories"`
	Count         int             `json:"count"`
}
