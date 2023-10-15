package model

import "realtime-chat/internal/otp"

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Status string  `json:"status,omitempty"`
	OTP    otp.OTP `json:"otp,omitempty"`
	Error  string  `json:"error,omitempty"`
}
