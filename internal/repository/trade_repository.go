package repository

import (
	"github.com/MonalBarse/tradelog/internal/domain"
	"gorm.io/gorm"
)

type TradeRepository interface {
	Create(trade *domain.Trade) error
	GetByUserID(userID uint) ([]domain.Trade, error)
	GetAll() ([]domain.Trade, error) // For Admins
}

type tradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository(db *gorm.DB) TradeRepository {
	return &tradeRepository{db}
}

func (r *tradeRepository) Create(trade *domain.Trade) error {
	return r.db.Create(trade).Error
}

// @desc: get trades for a specific user
func (r *tradeRepository) GetByUserID(userID uint) ([]domain.Trade, error) {
	var trades []domain.Trade
	// Preload is optional here, but good if we had complex relations
	err := r.db.Where("user_id = ?", userID).Find(&trades).Error
	return trades, err
}

// @desc: get all trades (admin)
func (r *tradeRepository) GetAll() ([]domain.Trade, error) {
	var trades []domain.Trade
	err := r.db.Find(&trades).Error
	return trades, err
}