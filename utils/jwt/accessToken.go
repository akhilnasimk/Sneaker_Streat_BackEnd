package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GenerateAccess(userId uuid.UUID, userName string, userEmail string, userRole string) (string, error) {
	secretcode := os.Getenv("Jwt_Secret")
	claim := jwt.MapClaims{
		"UserId":    userId.String(),
		"UserName":  userName,
		"UserEmail": userEmail,
		"UserRole":  userRole,
		"type":      "access",
		"exp":       time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwttocken, err := token.SignedString([]byte(secretcode))
	if err != nil {
		return "", err
	}
	return jwttocken, nil
}

func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("Jwt_Secret")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token type
		if claims["type"] != "access" {
			return nil, errors.New("invalid token type")
		}

		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return nil, errors.New("token expired")
			}
		} else {
			return nil, errors.New("invalid exp in token")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}
