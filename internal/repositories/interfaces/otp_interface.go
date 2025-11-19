package interfaces

import "github.com/akhilnasimk/SS_backend/internal/models"

type OtpRepository interface {
	SaveOtp(otp models.OTP) error
	FindOtpByEmailAndPurpose(email, purpose string) (*models.OTP, error)
}
