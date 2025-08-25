package middleware

import (
	"dekamond-task/package/jwt"
	"net/http"
	"strings"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"message":"missing token"}`))
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		_, err := jwt.ValidateJWT(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"message":"invalid token"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
