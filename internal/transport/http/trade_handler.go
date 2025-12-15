package http

import (
	"net/http"

	"github.com/MonalBarse/tradelog/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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
	Symbol   string          `json:"symbol" binding:"required"`
	Type     string          `json:"type" binding:"required,oneof=BUY SELL"` // restrict to BUY or SELL
	Price    decimal.Decimal `json:"price" binding:"required"`
	Quantity decimal.Decimal `json:"quantity" binding:"required"`
}

// Swagger Annotations
// @Summary Create a new trade
// @Description Records a buy or sell order. Validates sufficient funds for SELL orders.
// @Tags trades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createTradeRequest true "Trade Details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /trades [post]
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
	err := h.service.CreateTrade(c.Request.Context(), userID.(uint), req.Symbol, req.Type, req.Price, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Trade logged successfully"})
}

// Swagger Annotations
// @Summary List user trades
// @Description Get all trades for the logged-in user
// @Tags trades
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /trades [get]
func (h *TradeHandler) ListTrades(c *gin.Context) {
	userID, _ := c.Get("userID")

	trades, err := h.service.GetUserTrades(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trades})
}

// @Summary Get All Trades (Admin Only)
// @Description Get all trades across all users (admin only)
// @Tags trades
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /trades/all [get]
func (h *TradeHandler) GetAllTrades(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admins only"})
		return
	}

	trades, err := h.service.GetAllTrades(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all trades"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trades})
}

// @Summary Get Portfolio
// @Description Get current holdings calculated from trade history
// @Tags trades
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /portfolio [get]
func (h *TradeHandler) GetPortfolio(c *gin.Context) {
	userID, _ := c.Get("userID")

	portfolio, err := h.service.GetPortfolio(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate portfolio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": portfolio})
}
