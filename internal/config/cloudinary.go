package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

// CLD will hold the initialized Cloudinary instance
var CLD *cloudinary.Cloudinary

func InitCloudinary() {
	var err error
	CLD, err = cloudinary.NewFromParams(
		AppConfig.CLOUDINARY_CLOUD_NAME,
		AppConfig.CLOUDINARY_API_KEY,
		AppConfig.CLOUDINARY_API_SECRET,
	)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
}
