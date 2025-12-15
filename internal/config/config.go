package config

import (
	"log"

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

	if err := viper.ReadInConfig(); err != nil {
		color.Yellow("No .env file found, relying on environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		color.Red("Failed to unmarshal config: %v", err)
		log.Fatal("Could not load configuration")
	}

	if config.JWTSecret == "" || config.JWTRefreshSecret == "" {
		color.Red("JWT secrets are not set in configuration")
		log.Fatal("JWT secrets are required in configuration")
	}

	AppConfig = &config
}
