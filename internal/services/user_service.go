package services

import (
	"context"
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(id uuid.UUID) (dto.UserProfileResponse, error)
	UpdateProfile(userID uuid.UUID, profile *dto.UpdateProfileRequest) error
	GetAllUserData(limit, offset int) ([]dto.AdminUserResponse, error)
	GetUserById(ctx context.Context, stringID string) (dto.AdminUserResponse, error)
	AdminUserUpdate(ctx context.Context, req dto.PatchUserAdminReq, idstring string) error
}

type userService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) UserService {
	return &userService{
		userRepo: repo,
	}
}

func (S *userService) GetProfile(id uuid.UUID) (dto.UserProfileResponse, error) {
	var profile dto.UserProfileResponse //response

	user, err := S.userRepo.FindByID(id)

	if err != nil {
		return profile, err
	}

	if user == nil {
		return profile, fmt.Errorf("user not found")
	}

	// map model → DTO
	profile = dto.UserProfileResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Image:     user.Image,
		Phone:     user.Phone,
		Address:   user.Address,
		IsBlocked: user.IsBlocked,
		UserRole:  user.UserRole,
	}

	return profile, nil

}

func (s *userService) UpdateProfile(userID uuid.UUID, profile *dto.UpdateProfileRequest) error {
	//user exist or not
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	//collecting only value send by the user
	updates := make(map[string]interface{})

	if profile.UserName != nil && *profile.UserName != "" {
		updates["user_name"] = *profile.UserName
	}

	if profile.Email != nil && *profile.Email != "" {
		updates["email"] = *profile.Email
	}

	if profile.Image != nil {
		updates["image"] = *profile.Image
	}

	if profile.Phone != nil {
		updates["phone"] = *profile.Phone
	}

	if profile.Address != nil {
		updates["address"] = *profile.Address
	}

	// No fields to update
	if len(updates) == 0 {
		return fmt.Errorf("no valid fields provided to update")
	}

	// 3. Call repository
	return s.userRepo.PatchUser(userID, updates)
}

func (s *userService) GetAllUserData(limit, offset int) ([]dto.AdminUserResponse, error) {
	// Fetch paginated users from repository
	users, _, err := s.userRepo.GetAllUsersPaginated(limit, offset)
	if err != nil {
		return nil, err
	}

	// Map repository users to AdminUserResponse DTO
	var userResponses []dto.AdminUserResponse
	for _, u := range users {
		userResponses = append(userResponses, dto.AdminUserResponse{
			ID:        u.ID,
			UserName:  u.UserName,
			Email:     u.Email,
			Image:     u.Image,
			IsAdmin:   u.IsAdmin,
			IsBlocked: u.IsBlocked,
			CreatedAt: u.CreatedAt,
			UserRole:  u.UserRole,
		})
	}

	return userResponses, nil
}

func (s *userService) GetUserById(ctx context.Context, stringID string) (dto.AdminUserResponse, error) {
	// Convert string to UUID
	id := helpers.StringToUUID(stringID)
	if id == uuid.Nil {
		return dto.AdminUserResponse{}, fmt.Errorf("invalid UUID")
	}

	// Fetch user from repository
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return dto.AdminUserResponse{}, err
	}
	if user == nil {
		return dto.AdminUserResponse{}, fmt.Errorf("user not found")
	}

	// Map models.User -> dto.AdminUserResponse
	userResponse := dto.AdminUserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Image:     user.Image,
		IsAdmin:   user.IsAdmin,
		IsBlocked: user.IsBlocked,
		CreatedAt: user.CreatedAt,
		UserRole:  user.UserRole,
	}

	return userResponse, nil
}

func (s *userService) AdminUserUpdate(ctx context.Context, req dto.PatchUserAdminReq, idstring string) error {

	id := helpers.StringToUUID(idstring)
	if id == uuid.Nil {
		return fmt.Errorf("invalid user ID")
	}

	updates := make(map[string]interface{})
	if req.UserName != nil {
		updates["user_name"] = *req.UserName
	}
	if req.IsAdmin != nil {
		updates["is_admin"] = *req.IsAdmin
	}
	if req.IsBlocked != nil {
		updates["is_blocked"] = *req.IsBlocked
	}
	if req.UserRole != nil {
		updates["user_role"] = req.UserRole
	}
	if req.Image != nil {
		updates["image"] = req.Image
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	return s.userRepo.PatchUser(id, updates)
}
