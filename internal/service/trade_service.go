package service

import (
	"errors"
	"time"

	"github.com/MonalBarse/tradelog/internal/domain"
	"github.com/MonalBarse/tradelog/internal/repository"
)

// PortfolioItem represents the user's holding of a specific asset
type PortfolioItem struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Value    float64 `json:"value"` // Current market value (we'll just use last price for now)
}

type TradeService interface {
	CreateTrade(userID uint, symbol, tradeType string, price, quantity float64) error
	GetUserTrades(userID uint) ([]domain.Trade, error)
	GetAllTrades() ([]domain.Trade, error)
	GetPortfolio(userID uint) ([]PortfolioItem, error)
}

type tradeService struct {
	repo repository.TradeRepository
}

func NewTradeService(repo repository.TradeRepository) TradeService {
	return &tradeService{repo}
}

func (s *tradeService) CreateTrade(userID uint, symbol, tradeType string, price, quantity float64) error {
	// LOGIC 1: Validation
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	// LOGIC 2: Sell Constraints
	if tradeType == "SELL" {
		currentBalance, err := s.calculatePosition(userID, symbol)
		if err != nil {
			return err
		}

		if currentBalance < quantity {
			return errors.New("insufficient funds: you cannot sell more than you own")
		}
	}

	// Logic 3: Create the trade
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

// @desc: get portfolio for user
// @flow: get trades -> aggregate by symbol -> return holdings
func (s *tradeService) GetPortfolio(userID uint) ([]PortfolioItem, error) {
	trades, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Map to aggregate totals
	holdings := make(map[string]float64)

	for _, t := range trades {
		if t.Type == "BUY" {
			holdings[t.Symbol] += t.Quantity
		} else {
			holdings[t.Symbol] -= t.Quantity
		}
	}

	var portfolio []PortfolioItem
	for symbol, qty := range holdings {
		if qty > 0 { // Only show assets we actually own
			portfolio = append(portfolio, PortfolioItem{
				Symbol:   symbol,
				Quantity: qty,
				// Irl app, I would fetch live price here.
				// For now, I'll leave Value as 0 or calculate based on last trade.
			})
		}
	}

	return portfolio, nil
}
func (s *tradeService) calculatePosition(userID uint, symbol string) (float64, error) {
	trades, err := s.repo.GetByUserID(userID)
	if err != nil {
		return 0, err
	}

	var balance float64
	for _, t := range trades {
		if t.Symbol == symbol {
			if t.Type == "BUY" {
				balance += t.Quantity
			} else if t.Type == "SELL" {
				balance -= t.Quantity
			}
		}
	}
	return balance, nil
}
