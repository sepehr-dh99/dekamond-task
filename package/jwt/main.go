package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my-secret-key") // use env var in real apps

// CreateJWT generates a signed JWT for given phone.
func CreateJWT(phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": phone,
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	return token.SignedString(jwtKey)
}

// JWTAuth is middleware that checks Bearer token.
func JWTAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"missing token"}`, 401)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
