package service

import (
	"context"
	"errors"
	"time"

	"github.com/MonalBarse/tradelog/internal/domain"
	"github.com/MonalBarse/tradelog/internal/repository"
	"github.com/shopspring/decimal"
)

// PortfolioItem represents the user's holding of a specific asset
type PortfolioItem struct {
	Symbol   string          `json:"symbol"`
	Quantity decimal.Decimal `json:"quantity"`
	Value    decimal.Decimal `json:"value"` // Current market value (we'll just use last price for now)
}

type TradeService interface {
	CreateTrade(ctx context.Context, userID uint, symbol, tradeType string, price, quantity decimal.Decimal) error
	GetUserTrades(ctx context.Context, userID uint) ([]domain.Trade, error)
	GetAllTrades(ctx context.Context) ([]domain.Trade, error)
	GetPortfolio(ctx context.Context, userID uint) ([]PortfolioItem, error)
}

type tradeService struct {
	repo repository.TradeRepository
}

func NewTradeService(repo repository.TradeRepository) TradeService {
	return &tradeService{repo}
}

// @desc: create trade
// @flow: validate SELL -> check funds -> create trade record
func (s *tradeService) CreateTrade(ctx context.Context, userID uint, symbol, tradeType string, price, quantity decimal.Decimal) error {
	if quantity.LessThanOrEqual(decimal.Zero) { // quantity <= 0
		return errors.New("quantity must be positive")
	}

	if price.LessThanOrEqual(decimal.Zero) {
		return errors.New("price must be positive")
	}

	if tradeType == "SELL" {
		currentBalance, err := s.calculatePosition(ctx, userID, symbol)
		if err != nil {
			return err
		}

		if currentBalance.LessThan(quantity) { // currentBalance < quantity
			return errors.New("insufficient funds: you cannot sell more than you own")
		}
	}

	trade := &domain.Trade{
		UserID:     userID,
		Symbol:     symbol,
		Type:       tradeType,
		Price:      price,
		Quantity:   quantity,
		ExecutedAt: time.Now(),
	}
	return s.repo.Create(ctx, trade)
}

func (s *tradeService) GetUserTrades(ctx context.Context, userID uint) ([]domain.Trade, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *tradeService) GetAllTrades(ctx context.Context) ([]domain.Trade, error) {
	return s.repo.GetAll(ctx)
}

// @desc: get portfolio for user
// @flow: get trades -> aggregate by symbol -> return holdings
func (s *tradeService) GetPortfolio(ctx context.Context, userID uint) ([]PortfolioItem, error) {
	trades, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	holdings := make(map[string]decimal.Decimal)

	for _, t := range trades {
		if t.Type == "BUY" {
			holdings[t.Symbol] = holdings[t.Symbol].Add(t.Quantity)
		} else {
			holdings[t.Symbol] = holdings[t.Symbol].Sub(t.Quantity)
		}
	}

	var portfolio []PortfolioItem
	for symbol, qty := range holdings {
		if qty.GreaterThan(decimal.Zero) { // qty > 0
			portfolio = append(portfolio, PortfolioItem{
				Symbol:   symbol,
				Quantity: qty,
				Value:    decimal.Zero, // Placeholder
			})
		}
	}

	return portfolio, nil
}

func (s *tradeService) calculatePosition(ctx context.Context, userID uint, symbol string) (decimal.Decimal, error) {
	trades, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return decimal.Zero, err
	}

	balance := decimal.Zero // Initialize 0
	for _, t := range trades {
		if t.Symbol == symbol {
			if t.Type == "BUY" {
				balance = balance.Add(t.Quantity) // +
			} else if t.Type == "SELL" {
				balance = balance.Sub(t.Quantity) // -
			}
		}
	}
	return balance, nil
}
