package main

import (
	"net/http"
	"os"

	"github.com/MonalBarse/tradelog/docs"
	"github.com/MonalBarse/tradelog/internal/config"
	"github.com/MonalBarse/tradelog/internal/repository"
	"github.com/MonalBarse/tradelog/internal/service"
	transport "github.com/MonalBarse/tradelog/internal/transport/http"
	"github.com/MonalBarse/tradelog/internal/transport/middleware"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title TradeLog API
// @version 1.0
// @description This is a sample trading application backend.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@tradelog.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	if err := godotenv.Load(); err != nil {
		color.Yellow("No .env file found")
	}
	config.LoadConfig()

	config.ConnectDB()

	// @use : Access repo layer to interact with DB
	userRepo := repository.NewUserRepository(config.DB)
	tradeRepo := repository.NewTradeRepository(config.DB)

	// @use: Access service layer to use logic (for now auth only)
	authService := service.NewAuthService(userRepo)
	tradeService := service.NewTradeService(tradeRepo)

	// @use: Access transport layer to handle HTTP requests
	authHandler := transport.NewAuthHandler(authService)
	tradeHandler := transport.NewTradeHandler(tradeService)

	r := gin.Default() // here i am creating a gin router with default middleware (logger and recovery)

	api := r.Group("/api/v1") // I'll version it so if -> breaking changes -> can create new version
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// these are protected routes: only accessible with valid JWT
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/trades", tradeHandler.CreateTrade)
			protected.GET("/trades", tradeHandler.ListTrades)
			protected.GET("/portfolio", tradeHandler.GetPortfolio)
			protected.GET("/admin/trades", tradeHandler.GetAllTrades)
		}
	}

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
