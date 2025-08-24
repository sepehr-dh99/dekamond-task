package otp

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// OTPService manages OTP generation and verification.
type OTPService struct {
	mu        sync.Mutex
	codes     map[string]string    // phone -> otp code
	expiresAt map[string]time.Time // phone -> expiration time
}

func NewOTPService() *OTPService {
	return &OTPService{
		codes:     make(map[string]string),
		expiresAt: make(map[string]time.Time),
	}
}

// GenerateOTP creates and stores an OTP for the given phone.
func (o *OTPService) GenerateOTP(phone string) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	otp := o.generateSecureOTP()
	o.codes[phone] = otp
	o.expiresAt[phone] = time.Now().Add(2 * time.Minute)
	// Print to console (simulate sending SMS)
	log.Println("OTP for", phone, "is", otp)
	return otp, nil
}

// ValidateOTP checks if the provided OTP is correct and not expired.
func (o *OTPService) ValidateOTP(phone, code string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	exp, exists := o.expiresAt[phone]
	if !exists || time.Now().After(exp) {
		return errors.New("OTP expired or not found")
	}
	if o.codes[phone] != code {
		return errors.New("invalid OTP")
	}
	// Successful validation; remove OTP so it can't be reused
	delete(o.codes, phone)
	delete(o.expiresAt, phone)
	return nil
}

// generateSecureOTP returns a 4-digit random OTP (crypto/rand).
func (o *OTPService) generateSecureOTP() string {
	b := make([]byte, 2)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("Error generating OTP:", err)
	}
	return fmt.Sprintf("%04d", int(b[0])%10000)
}
