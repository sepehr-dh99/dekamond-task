package main

import (
	"log"
	"net/http"

	"dekamond-task/controller"
	"dekamond-task/package/jwt"
	"dekamond-task/package/otp"
	ratelimiter "dekamond-task/package/rateLimiter"
	"dekamond-task/service"
)

func main() {
	// Initialize in-memory stores and services
	userSvc := service.NewUserService()
	otpSvc := otp.NewOTPService()
	limiter := ratelimiter.NewRateLimiter()

	// Create HTTP handlers
	authCtrl := controller.NewAuthController(otpSvc, userSvc, limiter)
	userCtrl := controller.NewUserController(userSvc)

	// Setup routes (using std http for simplicity)
	http.HandleFunc("/auth/request-otp", authCtrl.RequestOTPHandler)
	http.HandleFunc("/auth/verify", authCtrl.VerifyOTPHandler)

	// Protected user routes
	http.Handle("/users", jwt.JWTAuth(userCtrl.ListUsersHandler))
	http.Handle("/users/", jwt.JWTAuth(userCtrl.GetUserHandler))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
