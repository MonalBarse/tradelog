package http

import (
	"net/http"

	"github.com/MonalBarse/tradelog/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

type promoteRequest struct {
	Secret string `json:"secret" binding:"required"`
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

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body registerRequest true "Registration Details"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Login godoc
// @Summary      User Login
// @Description  Authenticates user and returns an access token. Sets a refresh token in an HTTP-only cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body loginRequest true "Login Credentials"
// @Success      200  {object}  map[string]string "Returns access_token"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
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

// Refresh godoc
// @Summary      Refresh Access Token
// @Description  Uses the HttpOnly refresh_token cookie to issue a new access token
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string "Returns new access_token"
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	// Call service to handle refresh logic
	newAccessToken, newRefreshToken, err := h.service.Refresh(c.Request.Context(), refreshTokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

// Logout godoc
// @Summary      Logout User
// @Description  Logs out the user by clearing the refresh token cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// name, value, maxAge (-1 means delete), path, domain, secure, httpOnly
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// @Summary Promote user to Admin
// @Description Promotes the current user to admin if the correct secret is provided
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body promoteRequest true "Admin Secret"
// @Success 200 {object} map[string]string
// @Router /auth/promote [post]
func (h *AuthHandler) Promote(c *gin.Context) {
	var req promoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.service.PromoteToAdmin(c.Request.Context(), userID.(uint), req.Secret)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}