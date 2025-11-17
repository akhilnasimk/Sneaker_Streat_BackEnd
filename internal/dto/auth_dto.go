package dto

type RegisterRequest struct {
	UserName string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
	Address  string  `json:"address"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

type VerifyOTPReq struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}
type ChangePasswordReq struct {
	Email       string `json:"email"`
	NewPassword string `json:"newpassword"`
}
