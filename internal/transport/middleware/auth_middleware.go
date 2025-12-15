package middleware

import (
	"net/http"
	"strings"

	"github.com/MonalBarse/tradelog/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// @desc : auth-middleware -> checks for valid jwt in auth header
// @flow : get auth header -> parse token -> validate -> extract claims -> inject into context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get auth header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// validate token
		token, err := utils.ValidateAccessToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// extract the payload
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// inject into gin context to use later in handlers
		// (We can retrieve this later in the handlers using c.Get("userID"))
		c.Set("userID", uint(claims["sub"].(float64)))
		c.Set("role", claims["role"])

		c.Next()
	}
}
