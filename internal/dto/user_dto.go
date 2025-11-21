package dto

import (
	"time"

	"github.com/google/uuid"
)

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

type AdminUserResponse struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Image     *string   `json:"image,omitempty"`
	IsAdmin   bool      `json:"is_admin"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
	UserRole  *string   `json:"role"`
}

type PatchUserAdminReq struct {
	UserName *string `json:"username,omitempty"`
	IsAdmin  *bool   `json:"is_admin,omitempty"`
	UserRole *string `json:"role,omitempty"`
	Image    *string `json:"image,omitempty"`
}
