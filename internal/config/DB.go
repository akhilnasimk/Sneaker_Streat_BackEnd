package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // ⬅️ NEW IMPORT: To load the .env file
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// 1. Load the .env file (CRITICAL STEP)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file. Using system environment variables.")
		// Note: You can change this to log.Fatal if you require the .env file to be present.
	}

	// 2. Construct the DSN string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)

	// 3. Open connection with GORM Config
	// NOTE: We wrap postgres.Open with postgres.New and Config{} to add the PGBouncer flag
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		// 🚨 CRITICAL FOR SUPABASE POOLER: Disables prepared statement caching
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	DB = db
	fmt.Println("🚀 Connected to Supabase Postgres!")
}
