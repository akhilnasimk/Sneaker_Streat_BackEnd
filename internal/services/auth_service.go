package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories"
	"github.com/akhilnasimk/SS_backend/utils/jwt"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(User dto.RegisterRequest) error
	Login(user dto.LoginReq) (*models.User, error)
	GenerateAndStoreRefreshToken(userID uuid.UUID) (string, error)
	RefreshTokens(token string) (string, string, error)
	ForgotPassword(email string) (*models.User, error)
	UpdatePassword(email string, newPassword string) error
}

type authService struct {
	userRepo  repositories.UserRepository
	JwtSecret []byte
	tokenRepo repositories.TokenRepository
}

func NewAuthService(repo repositories.UserRepository, tokenRepo repositories.TokenRepository) AuthService {
	return &authService{
		userRepo:  repo,
		tokenRepo: tokenRepo,
	}
}

func (S *authService) Register(User dto.RegisterRequest) error {
	user := models.User{
		UserName: User.UserName,
		Email:    User.Email,
		Password: User.Password,
		Phone:    &User.Password,
		Address:  &User.Address,
	}
	// checking if user exist
	existingUser, err := S.userRepo.FindByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Hashing pass
	hashedPass, err := helpers.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPass) // setting password

	// Registering user with the injected user repo methode
	if err := S.userRepo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}
	log.Println("regist sucess")
	return nil
}

func (S *authService) Login(req dto.LoginReq) (*models.User, error) {
	user, err := S.userRepo.FindByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	fmt.Println(user.Password)

	if !helpers.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (S *authService) GenerateAndStoreRefreshToken(userID uuid.UUID) (string, error) {
	//genarating
	refreshToken, _ := jwt.GenerateRefreshToken()

	// Hashing
	hashedToken := jwt.HashRefresh(refreshToken)

	// seting up the model to save
	token := models.RefreshToken{
		UserID:    userID,
		Token:     hashedToken,                        // store the hashed version only
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days validity
	}

	// pass the non hashed refresh to user and and cookie store hashed version to db
	return refreshToken, S.tokenRepo.SaveRefreshToken(&token)
}
func (S *authService) RefreshTokens(token string) (string, string, error) {

	hashedToken := jwt.HashRefresh(token)
	tokenResp, err := S.tokenRepo.FindByToken(hashedToken)
	if err != nil {
		return "", "", err
	}

	// 🔥 If no token record exists, return error
	if tokenResp == nil {
		return "", "", errors.New("refresh token not found or invalid")
	}

	if tokenResp.RevokedAt != nil {
		return "", "", errors.New("refresh token is revoked")
	}

	// Check expiry
	if time.Now().After(tokenResp.ExpiresAt) {
		return "", "", errors.New("refresh token has expired")
	}

	// Fetch user
	user, err := S.userRepo.FindByID(tokenResp.UserID)
	if err != nil {
		return "", "", err
	}

	// Create new access
	newAccess, err := jwt.GenerateAccess(user.ID, user.UserName, user.Email, *user.UserRole)
	if err != nil {
		return "", "", err
	}

	// Create a new refresh
	newRefresh, err := jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Hash it before saving
	newHashedToken := jwt.HashRefresh(newRefresh)

	if err := S.tokenRepo.UpdateToken(user.ID, newHashedToken, time.Now().Add(7*24*time.Hour)); err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}


func (S *authService) ForgotPassword(email string) (*models.User, error) {

	user, err := S.userRepo.FindByEmail(email)
	if err != nil {
		return &models.User{}, err
	}

	if user == nil {
		return &models.User{}, fmt.Errorf("user is Not found ")
	}
	if user.IsBlocked {
		return &models.User{}, fmt.Errorf("user is blocked")
	}
	return user, nil
}

func (S *authService) UpdatePassword(email string, newPassword string) error {
	if email == "" || newPassword == "" {
		return fmt.Errorf("no email or password has been send ")
	}

	hashed, err := helpers.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := S.userRepo.PatchPasswordByEmail(email, hashed); err != nil {
		return err
	}

	return nil

}
