package dto

import "github.com/google/uuid"

type UserProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Image     *string   `json:"image,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Address   *string   `json:"address,omitempty"`
	IsBlocked bool      `json:"is_blocked"`
	UserRole  *string   `json:"user_role"`
}

type UpdateProfileRequest struct {
	UserName *string `json:"username"`
	Email    *string `json:"email"`
	Image    *string `json:"image"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
}
