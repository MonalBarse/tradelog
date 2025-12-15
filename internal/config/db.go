package config

import (
	"log"

	"github.com/MonalBarse/tradelog/internal/domain"
	"github.com/fatih/color"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// @Desc: Esttablishes conn to PosgresDB (using GORM) and runs AutoMigrate

func ConnectDB() {
	// building data source name -> using env variables
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("DB_PORT"),
	// )
	dsn := AppConfig.DBUrl //shorteer way

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		color.Red("Failed to connect to DB: %v", err)
		log.Fatal("Failed to connect to DB:", err)
	}

	log.Println("db connected XOXO")

	//AutoMigrate (create tables automatically based on structs)
	// In production, I would use proper migration files, but for this assignment, AutoMigrate should be acceptable
	log.Println("Running migrations...")
	err = DB.AutoMigrate(&domain.User{}, &domain.Trade{})
	if err != nil {
		color.Red("Migration failed :( : %v", err)
		log.Fatal("Migration failed :(  :", err)
	}
	color.Green("----------------DB migrations done XOXO-----------------")
}
