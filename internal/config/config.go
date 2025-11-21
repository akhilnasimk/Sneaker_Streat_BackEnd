package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database config
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBPort string

	// JWT config
	JwtSecret string

	// SMTP/email config
	SMTPEmail             string
	SMTPPass              string
	SMTPHost              string
	SMTPPort              string
	CLOUDINARY_CLOUD_NAME string
	CLOUDINARY_API_KEY    string
	CLOUDINARY_API_SECRET string
}

// Global variable to hold the loaded config
var AppConfig *Config

func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		DBHost:                os.Getenv("DB_HOST"),
		DBUser:                os.Getenv("DB_USER"),
		DBPass:                os.Getenv("DB_PASS"),
		DBName:                os.Getenv("DB_NAME"),
		DBPort:                os.Getenv("DB_PORT"),
		JwtSecret:             os.Getenv("Jwt_Secret"),
		SMTPEmail:             os.Getenv("SMTP_EMAIL"),
		SMTPPass:              os.Getenv("SMTP_PASSWORD"),
		SMTPHost:              os.Getenv("SMTP_HOST"),
		SMTPPort:              os.Getenv("SMTP_PORT"),
		CLOUDINARY_CLOUD_NAME: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CLOUDINARY_API_KEY:    os.Getenv("CLOUDINARY_API_KEY"),
		CLOUDINARY_API_SECRET: os.Getenv("CLOUDINARY_API_SECRET"),
	}
}
