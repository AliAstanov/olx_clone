package models

import (
	"time"

	"github.com/google/uuid"
)

// // ProductStatus turi
type ProductStatus string

// Mavjud statuslar
const (
	StatusPending  ProductStatus = "pending"
	StatusActive   ProductStatus = "active"
	StatusRejected ProductStatus = "rejected"
	StatusExpires  ProductStatus = "expires"
)

// type ProductStatus struct {
// 	StatusPending  string
// 	StatusActive   string
// 	StatusRejected string
// 	StatusExpires  string
// }

// Product struct
type Product struct {
	ID            uuid.UUID     `json:"id"`
	UserID        uuid.UUID     `json:"user_id"` // O'zgartirildi
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	Price         float64       `json:"price"`
	Status        ProductStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	SubcategoryID uuid.UUID     `json:"subcategory_id"` // O'zgartirildi
	Image         string        `json:"image"`
	IsNew         bool          `json:"is_new"`
	ExpiresAt     *time.Time     `json:"expires_at"`
	Images        []string    `json:"images"` 
}

type CreateProductReq struct {
	UserID        uuid.UUID `json:"user_id" validate:"required"`
	SubcategoryID uuid.UUID `json:"subcategory_id" validate:"required"`
	Title         string    `json:"title" validate:"required,min=3,max=255"`
	Description   string    `json:"description" validate:"omitempty,max=5000"`
	Price         float64   `json:"price" validate:"required,gt=0"`
	Image         string    `json:"image" validate:"omitempty,url"`
	IsNew         bool      `json:"is_new"`
	ExpiresAt     time.Time `json:"expires_at" validate:"required"`
}

type UpdateProductReq struct {
	Title       string    `json:"title,omitempty" validate:"max=255"` // Ixtiyoriy, 255 belgigacha
	Description string    `json:"description,omitempty"`              // Ixtiyoriy
	Price       float64   `json:"price,omitempty" validate:"gt=0"`    // Ixtiyoriy, 0 dan katta bo'lishi kerak
	Image       string    `json:"image,omitempty"`                    // Ixtiyoriy
	IsNew       bool      `json:"is_new,omitempty"`                   // Ixtiyoriy, yangi yoki foydalanilganligini bildiradi
	ExpiresAt   time.Time `json:"expires_at,omitempty"`               // Ixtiyoriy, amal qilish muddati
}

type GetListProducts struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}

type ProductImage struct {
    ID        uuid.UUID `json:"id"`
    ProductID uuid.UUID `json:"product_id"`
    ImageURL  string    `json:"image_url"`
    CreatedAt time.Time `json:"created_at"`
}
