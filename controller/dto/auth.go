package dto

type RequestOTPRequest struct {
	Phone string `json:"phone" example:"09123456789" validate:"required,startswith=09,len=11"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" example:"09123456789" validate:"required,startswith=09,len=11"`
	OTP   string `json:"otp" example:"1234" validate:"required,len=6,numeric"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
}
