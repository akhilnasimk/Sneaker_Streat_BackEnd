package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	// Convert to 6-digit padded string
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// HashOTP hashes the OTP before saving to DB
func HashOTP(otp string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyOTP compares input OTP with the hashed OTP from DB
func VerifyOTPHash(input, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
	return err == nil
}
