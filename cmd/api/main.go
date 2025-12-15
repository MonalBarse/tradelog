package main

import (
	"net/http"
	"os"

	"github.com/MonalBarse/tradelog/internal/config"
	"github.com/MonalBarse/tradelog/internal/repository"
	"github.com/MonalBarse/tradelog/internal/service"
	transport "github.com/MonalBarse/tradelog/internal/transport/http"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		color.Yellow("No .env file found")
	}

	config.ConnectDB()

	// @use : Access repo layer to interact with DB
	userRepo := repository.NewUserRepository(config.DB)

	// @use: Access service layer to use logic (for now auth only)
	authService := service.NewAuthService(userRepo)

	// @use: Access transport layer to handle HTTP requests
	authHandler := transport.NewAuthHandler(authService)
	r := gin.Default() // here i am creating a gin router with default middleware (logger and recovery)

	api := r.Group("/api/v1") // I'll version it so if -> breaking changes -> can create new version
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     "connected",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	color.Cyan("--------------Server running on port %s ----------------\n ", port)
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
