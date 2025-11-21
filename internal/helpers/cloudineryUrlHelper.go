package helpers

import (
	"path/filepath"
	"strings"
)

func ExtractPublicID(url string) string {
	// URL example: https://res.cloudinary.com/demo/image/upload/v1234567890/products/shoe1.jpg
	parts := strings.Split(url, "/upload/")
	if len(parts) != 2 {
		return "" // invalid URL
	}
	// remove the version number and file extension
	publicWithVersion := parts[1]                                            // v1234567890/products/shoe1.jpg
	publicID := strings.Join(strings.Split(publicWithVersion, "/")[1:], "/") // remove version if needed
	publicID = strings.TrimSuffix(publicID, filepath.Ext(publicID))          // remove .jpg
	return publicID
}
