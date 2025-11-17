package services

import (
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(id uuid.UUID) (dto.UserProfileResponse, error)
	UpdateProfile(userID uuid.UUID, profile *dto.UpdateProfileRequest) error
	GetAllUserData(limit, offset int) ([]dto.AdminUserResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
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
