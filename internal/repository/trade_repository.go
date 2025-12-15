package repository

import (
	"context"

	"github.com/MonalBarse/tradelog/internal/domain"
	"gorm.io/gorm"
)

type TradeRepository interface {
	Create(ctx context.Context, trade *domain.Trade) error
	GetByUserID(ctx context.Context, userID uint) ([]domain.Trade, error)
	GetAll(ctx context.Context) ([]domain.Trade, error) // For Admins
}

type tradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository(db *gorm.DB) TradeRepository {
	return &tradeRepository{db}
}

func (r *tradeRepository) Create(ctx context.Context, trade *domain.Trade) error {
	return r.db.WithContext(ctx).Create(trade).Error
}

// @desc: get trades for a specific user
func (r *tradeRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Trade, error) {
	var trades []domain.Trade
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&trades).Error
	return trades, err
}

// @desc: get all trades (admin)
func (r *tradeRepository) GetAll(ctx context.Context) ([]domain.Trade, error) {
	var trades []domain.Trade
	err := r.db.WithContext(ctx).Find(&trades).Error
	return trades, err
}
