package http

import (
	"net/http"

	"github.com/MonalBarse/tradelog/internal/service"
	"github.com/gin-gonic/gin"
)

/*
NOTE:
	I use:
		c.Get("userID"): This retrieves the data I stored in the AuthMiddleware
		Trust: I trust this ID because the middleware already validated the token
*/

type TradeHandler struct {
	service service.TradeService
}

func NewTradeHandler(service service.TradeService) *TradeHandler {
	return &TradeHandler{service}
}

type createTradeRequest struct {
	Symbol   string  `json:"symbol" binding:"required"`
	Type     string  `json:"type" binding:"required,oneof=BUY SELL"` // restrict to BUY or SELL
	Price    float64 `json:"price" binding:"required,gt=0"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

// @desc - create a new trade for logged-in user
// @req  - POST /trades
// @flow - bind json -> get userID from context -> call service -> res
func (h *TradeHandler) CreateTrade(c *gin.Context) {
	var req createTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get uid from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// actualy create trade
	err := h.service.CreateTrade(userID.(uint), req.Symbol, req.Type, req.Price, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trade"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Trade logged successfully"})
}

// @desc - list trades for logged-in user
// @req  - get /trades
// @flow - get userID from context -> call service -> res
func (h *TradeHandler) ListTrades(c *gin.Context) {
	userID, _ := c.Get("userID")

	trades, err := h.service.GetUserTrades(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trades})
}

// @desc - admin: get all trades
// @req  - get /admin/trades
// @flow - check role from context -> call service -> res
func (h *TradeHandler) GetAllTrades(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admins only"})
		return
	}

	trades, err := h.service.GetAllTrades()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all trades"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trades})
}
