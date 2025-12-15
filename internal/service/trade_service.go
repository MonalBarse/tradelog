package service

import (
	"time"

	"github.com/MonalBarse/tradelog/internal/domain"
	"github.com/MonalBarse/tradelog/internal/repository"
)

type TradeService interface {
	CreateTrade(userID uint, symbol, tradeType string, price, quantity float64) error
	GetUserTrades(userID uint) ([]domain.Trade, error)
	GetAllTrades() ([]domain.Trade, error)
}

type tradeService struct {
	repo repository.TradeRepository
}

func NewTradeService(repo repository.TradeRepository) TradeService {
	return &tradeService{repo}
}

// @desc: create a new trade
func (s *tradeService) CreateTrade(userID uint, symbol, tradeType string, price, quantity float64) error {
	trade := &domain.Trade{
		UserID:     userID,
		Symbol:     symbol,
		Type:       tradeType,
		Price:      price,
		Quantity:   quantity,
		ExecutedAt: time.Now(),
	}
	return s.repo.Create(trade)
}

func (s *tradeService) GetUserTrades(userID uint) ([]domain.Trade, error) {
	return s.repo.GetByUserID(userID)
}

func (s *tradeService) GetAllTrades() ([]domain.Trade, error) {
	return s.repo.GetAll()
}
