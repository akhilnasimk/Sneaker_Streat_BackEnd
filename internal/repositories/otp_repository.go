package repositories

import (
	"errors"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OtpRepository interface {
	SaveOtp(otp models.OTP) error
	FindOtpByEmailAndPurpose(email, purpose string) (*models.OTP, error)
}

type otpRepository struct {
	DB *gorm.DB
}

func NewOtpRepository(db *gorm.DB) OtpRepository {
	return &otpRepository{
		DB: db,
	}
}

func (r otpRepository) SaveOtp(otp models.OTP) error {
	// Upsert: insert new OTP or update existing one if email + purpose conflict
	err := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}, {Name: "purpose"}},
		DoUpdates: clause.AssignmentColumns([]string{"otp_code", "expires_at", "is_used", "user_id"}),
	}).Create(&otp).Error

	return err
}

// FindOtpByEmailAndPurpose fetches the OTP record for a given email and purpose
func (r otpRepository) FindOtpByEmailAndPurpose(email, purpose string) (*models.OTP, error) {
	var otp models.OTP
	err := r.DB.Where("email = ? AND purpose = ?", email, purpose).
		Order("created_at desc").
		First(&otp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // not found
		}
		return nil, err
	}

	return &otp, nil
}
