package http

import (
	"net/http"

	"github.com/MonalBarse/tradelog/internal/service"
	"github.com/MonalBarse/tradelog/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

// Request Data Structures for Binding JSON
type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @desc Creates a new user account with email and password
// @req  post /auth/register
// @flow  validate input -> call register user -> res
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Register(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// @desc Authenticates user and returns JWT token
// @req post /auth/login
// @flow validate input -> call login user -> set cookie + res
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// SET HTTP-ONLY COOKIE
	// name, value, maxAge (seconds), path, domain, secure, httpOnly
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true) // Secure=false for localhost. Set true in prod.

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

// @desc Refreshes JWT tokens using the refresh token cookie
// @req  post /auth/refresh
// @flow get cookie -> validate -> generate new tokens -> set cookie + res
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	token, err := utils.ValidateRefreshToken(refreshTokenString)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID := uint(claims["sub"].(float64))

	// In a real app, we might fetch the user role from DB here to ensure they aren't banned.
	// For now, we assume "user" or store role in refresh token too.
	newAccessToken, newRefreshToken, err := utils.GenerateTokens(userID, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
