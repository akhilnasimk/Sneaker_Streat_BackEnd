package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/akhilnasimk/SS_backend/utils/otp"
	"github.com/google/uuid"
)

type OtpService interface {
	SendOTP(userID uuid.UUID, email string, purpose string) error
	VerifyOTP(inputOtp string, email string, purpose string) (bool, error)
}

type otpService struct {
	otpRepo      interfaces.OtpRepository
	EmailService EmailService
}

func NewOtpService(repo interfaces.OtpRepository, ES EmailService) OtpService {
	return &otpService{
		otpRepo:      repo,
		EmailService: ES,
	}
}

func (R *otpService) SendOTP(userID uuid.UUID, email string, purpose string) error {
	// generate OTP
	otpstring, err := otp.GenerateOTP()
	if err != nil {
		return err
	}

	//hash otp
	hashedOtp, _ := otp.HashOTP(otpstring)
	
	// store OTP FIRST (so user can verify even if email fails)
	if err := R.otpRepo.SaveOtp(models.OTP{
		UserID:    &userID,
		OTPCode:   hashedOtp,
		Email:     email,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Purpose:   purpose,
	}); err != nil {
		return err
	}

	// Send email asynchronously (don't wait for it)
	go func() {
		if err := R.EmailService.SendMailOTP(email, otpstring); err != nil {
			// Log the error but don't fail the request
			fmt.Printf("Failed to send OTP email to %s: %v\n", email, err)
		}
	}()

	return nil
}

func (R *otpService) VerifyOTP(inputOtp string, email string, purpose string) (bool, error) {
	otpResp, err := R.otpRepo.FindOtpByEmailAndPurpose(email, purpose)
	if err != nil { //Fail check
		return false, err
	}
	//Nil chcekc
	if otpResp == nil {
		return false, errors.New("OTP not found")
	}

	//checking if the otp is  usesed
	if otpResp.IsUsed {
		return false, errors.New("OTP already used")
	}

	//check the expiry
	if time.Now().After(otpResp.ExpiresAt) {
		return false, errors.New("OTP expired")
	}

	// Compare the OTP
	if !otp.VerifyOTPHash(inputOtp, otpResp.OTPCode) {
		return false, errors.New("invalid OTP")
	}

	// Mark as used
	otpResp.IsUsed = true
	if err := R.otpRepo.SaveOtp(*otpResp); err != nil {
		return false, err
	}

	return true, nil
}
