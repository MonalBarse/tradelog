package config

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Config struct {
	Env              string `mapstructure:"ENV"`
	Port             string `mapstructure:"PORT"`
	DBUrl            string `mapstructure:"DB_URL"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
	JWTRefreshSecret string `mapstructure:"JWT_REFRESH_SECRET"`
	AdminSecret      string `mapstructure:"ADMIN_SECRET"`
}

var AppConfig *Config // Global accessible config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Try to read from .env file, but don't fail if it doesn't exist
	if err := viper.ReadInConfig(); err != nil {
		color.Yellow("No .env file found, relying on environment variables")
	}

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		color.Red("Failed to unmarshal config: %v", err)
		log.Fatal("Could not load configuration")
	}

	if config.JWTSecret == "" {
		config.JWTSecret = os.Getenv("JWT_SECRET")
	}
	if config.JWTRefreshSecret == "" {
		config.JWTRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")
	}
	if config.AdminSecret == "" {
		config.AdminSecret = os.Getenv("ADMIN_SECRET")
	}
	if config.DBUrl == "" {
		config.DBUrl = os.Getenv("DB_URL")
	}
	if config.Port == "" {
		config.Port = os.Getenv("PORT")
		if config.Port == "" {
			config.Port = "8080"
		}
	}
	if config.Env == "" {
		config.Env = os.Getenv("ENV")
		if config.Env == "" {
			config.Env = "development"
		}
	}

	// Validate required secrets
	if config.JWTSecret == "" || config.JWTRefreshSecret == "" {
		color.Red("JWT secrets are not set in configuration")
		color.Red("JWT_SECRET: %s", config.JWTSecret)
		color.Red("JWT_REFRESH_SECRET: %s", config.JWTRefreshSecret)
		log.Fatal("JWT secrets are required in configuration")
	}

	color.Green("Configuration loaded successfully!")
	color.Cyan("Port: %s", config.Port)
	color.Cyan("Environment: %s", config.Env)
	color.Cyan("DB URL: %s", config.DBUrl)

	AppConfig = &config
}
