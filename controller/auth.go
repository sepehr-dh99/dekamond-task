package controller

import (
	"encoding/json"
	"net/http"

	"dekamond-task/controller/dto"
	"dekamond-task/package/jwt"
	"dekamond-task/package/otp"
	ratelimiter "dekamond-task/package/rate_limiter"
	"dekamond-task/package/response"
	"dekamond-task/package/validator"
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
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RequestOTPRequest true "Phone number (09XXXXXXXXX)"
// @Success 200 {object} response.Response[any] "Successful operation"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 429 {object} response.ErrorResponse "Too many requests"
// @Router /auth/request-otp [post]
func (ac *AuthController) RequestOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestOTPRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Validate DTO
	if err := validator.Validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "phone is required")
		return
	}

	// Rate limit check
	if err := ac.limiter.Allow(req.Phone); err != nil {
		response.Error(w, http.StatusTooManyRequests, "too many requests")
		return
	}

	// Generate and store OTP
	ac.otpSvc.GenerateOTP(req.Phone)
	response.Success[any](w, nil, "OTP sent successfully")
}

// VerifyOTPHandler handles POST /auth/verify.
// @Summary Verify OTP and login/register
// @Description Validates OTP, registers user if new, and returns JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Phone and OTP"
// @Success 200 {object} response.Response[dto.VerifyOTPResponse] "Login successful"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/verify [post]
func (ac *AuthController) VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyOTPRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Validate DTO
	if err := validator.Validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "phone and otp required")
		return
	}

	// Validate OTP
	if err := ac.otpSvc.ValidateOTP(req.Phone, req.OTP); err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Register or fetch existing user
	user := ac.userSvc.RegisterIfNotExists(req.Phone)

	// Generate JWT
	token, _ := jwt.CreateJWT(user.Phone)
	response.Success(w, &dto.VerifyOTPResponse{Token: token}, "login successful")
}
