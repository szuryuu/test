package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UmkmID          uuid.UUID  `json:"umkm_id" db:"umkm_id"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty" db:"category_id"`
	Amount          int64      `json:"amount" db:"amount"`
	Type            string     `json:"type" db:"type"` // income | expense
	Description     string     `json:"description" db:"description"`
	RawMessage      string     `json:"raw_message" db:"raw_message"`
	TransactionDate string     `json:"transaction_date" db:"transaction_date"` // YYYY-MM-DD
	Source          string     `json:"source" db:"source"`                     // whatsapp | manual
	AIConfidence    *float64   `json:"ai_confidence,omitempty" db:"ai_confidence"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type TransactionCategory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	Icon      *string   `json:"icon,omitempty" db:"icon"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
