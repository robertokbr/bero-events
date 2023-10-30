package middlewares

import (
	"fmt"
	"net/http"
	"os"
)

func RequestAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authHeader := r.Header.Get("Authorization")
	secret := os.Getenv("BASIC_AUTH_VALUE")
	authorization := fmt.Sprintf("Basic %s", secret)

	if authHeader != authorization {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	next(w, r)
}
