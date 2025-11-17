package migrations

import (
	"log"

	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/models"
)

func RunMigrations() {
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Category{},
		&models.Order{},
		&models.OrderItem{},
		&models.ProductImage{},
		&models.Wishlist{},
		&models.RefreshToken{},
		&models.OTP{},
	)
	if err != nil {
		log.Fatal("Migration failed ", err)
	}
}
