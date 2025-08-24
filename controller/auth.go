package controller

import (
	"encoding/json"
	"net/http"

	"dekamond-task/package/jwt"
	"dekamond-task/package/otp"
	ratelimiter "dekamond-task/package/rateLimiter"
	"dekamond-task/service"
)

type AuthController struct {
	otpSvc  *otp.OTPService
	userSvc *service.UserService
	limiter *ratelimiter.RateLimiter
}

func NewAuthController(o *otp.OTPService, u *service.UserService, l *ratelimiter.RateLimiter) *AuthController {
	return &AuthController{otpSvc: o, userSvc: u, limiter: l}
}

// RequestOTPHandler handles POST /auth/request-otp.
// @Summary Request OTP
// @Description Generate OTP for the given phone (Iranian format).
// @Param phone body string true "Phone number (09XXXXXXXXX)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /auth/request-otp [post]
func (ac *AuthController) RequestOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	phone := req.Phone
	if phone == "" {
		http.Error(w, `{"error":"phone is required"}`, 400)
		return
	}
	// Rate limit check
	if err := ac.limiter.Allow(phone); err != nil {
		http.Error(w, `{"error":"too many requests"}`, 429)
		return
	}
	// Generate and store OTP
	ac.otpSvc.GenerateOTP(phone)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent successfully"})
}

// VerifyOTPHandler handles POST /auth/verify.
// @Summary Verify OTP and login/register
// @Description Validates OTP, registers user if new, and returns JWT.
// @Param phone body string true "Phone number"
// @Param otp   body string true "OTP code"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/verify [post]
func (ac *AuthController) VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Phone == "" || req.OTP == "" {
		http.Error(w, `{"error":"phone and otp required"}`, 400)
		return
	}
	// Validate OTP
	if err := ac.otpSvc.ValidateOTP(req.Phone, req.OTP); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, 401)
		return
	}
	// Register or fetch existing user
	user := ac.userSvc.RegisterIfNotExists(req.Phone)
	// Generate JWT
	token, _ := jwt.CreateJWT(user.Phone) // see middleware/jwt_middleware.go
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
