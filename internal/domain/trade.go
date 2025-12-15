package domain

import (
	"time"

	"gorm.io/gorm"
)

type Trade struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"` // Foreign Key with Index
	Symbol      string         `gorm:"not null" json:"symbol"`        // e.g., "BTC/USD"
	Type        string         `gorm:"not null" json:"type"`          // "BUY" or "SELL"
	Price       float64        `gorm:"not null;type:decimal(20,8)" json:"price"` // High precision for crypto
	Quantity    float64        `gorm:"not null;type:decimal(20,8)" json:"quantity"`
	Notes       string         `json:"notes"`
	ExecutedAt  time.Time      `gorm:"not null" json:"executed_at"` // When the trade happened
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}