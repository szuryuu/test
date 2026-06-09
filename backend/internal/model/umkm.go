package model

import (
	"time"

	"github.com/google/uuid"
)

type Umkm struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	BusinessName string     `json:"business_name" db:"business_name"`
	PhoneNumber  string     `json:"phone_number" db:"phone_number"`
	Email        *string    `json:"email,omitempty" db:"email"`
	Address      *string    `json:"address,omitempty" db:"address"`
	BusinessType *string    `json:"business_type,omitempty" db:"business_type"`
	PasswordHash string     `json:"-" db:"password_hash"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
