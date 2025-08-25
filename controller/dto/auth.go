package dto

type RequestOTPRequest struct {
	Phone string `json:"phone" example:"09123456789"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" example:"09123456789"`
	OTP   string `json:"otp" example:"1234"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
}
