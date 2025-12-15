package service

import (
	"context"
	"testing"

	"github.com/MonalBarse/tradelog/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockTradeRepo struct {
	mock.Mock
}

func (m *MockTradeRepo) Create(ctx context.Context, trade *domain.Trade) error {
	args := m.Called(ctx, trade)
	return args.Error(0)
}

func (m *MockTradeRepo) GetByUserID(ctx context.Context, userID uint) ([]domain.Trade, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Trade), args.Error(1)
}

func (m *MockTradeRepo) GetAll(ctx context.Context) ([]domain.Trade, error) {
	return nil, nil
}

func TestCreateTrade_InsufficientFunds(t *testing.T) {
	// Setup
	mockRepo := new(MockTradeRepo)
	service := NewTradeService(mockRepo)
	ctx := context.Background()

	// Mock: User has bought 10 BTC previously
	mockRepo.On("GetByUserID", ctx, uint(1)).Return([]domain.Trade{
		{Symbol: "BTC/USD", Type: "BUY", Quantity: decimal.NewFromInt(10)},
	}, nil)

	// Attempt to Sell 20 BTC
	err := service.CreateTrade(ctx, 1, "BTC/USD", "SELL", decimal.NewFromInt(50000), decimal.NewFromInt(20))

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds: you cannot sell more than you own", err.Error())
	mockRepo.AssertExpectations(t)
}
