package main

import (
	"log"
	"net/http"

	"dekamond-task/controller"
	"dekamond-task/package/jwt"
	"dekamond-task/package/otp"
	ratelimiter "dekamond-task/package/rateLimiter"
	"dekamond-task/service"

	_ "dekamond-task/docs"

	httpSwagger "github.com/swaggo/http-swagger" // Swagger UI handler
)

// @title           Dekamond Task API
// @version         1.0
// @description     OTP-based authentication and user management API.
// @BasePath        /
// @schemes         http
// @host            localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize in-memory stores and services
	userSvc := service.NewUserService()
	otpSvc := otp.NewOTPService()
	limiter := ratelimiter.NewRateLimiter()

	// Create HTTP handlers
	authCtrl := controller.NewAuthController(otpSvc, userSvc, limiter)
	userCtrl := controller.NewUserController(userSvc)

	// Public auth routes
	http.HandleFunc("/auth/request-otp", authCtrl.RequestOTPHandler)
	http.HandleFunc("/auth/verify", authCtrl.VerifyOTPHandler)

	// Protected user routes
	http.Handle("/users", jwt.JWTAuth(userCtrl.ListUsersHandler))
	http.Handle("/users/", jwt.JWTAuth(userCtrl.GetUserHandler))

	// Swagger UI (visit http://localhost:8080/swagger/index.html)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
