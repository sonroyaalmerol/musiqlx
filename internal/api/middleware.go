package api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/sonroyaalmerol/musiqlx/internal/services"
)

func ProtectedPath(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return AuthMiddleware(http.HandlerFunc(next))
}

func UnprotectedPath(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(next)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract query parameters
		username := r.URL.Query().Get("u")
		password := r.URL.Query().Get("p")
		version := r.URL.Query().Get("v")
		salt := r.URL.Query().Get("s")
		token := r.URL.Query().Get("t")

		// Validate API version and perform corresponding authentication
		if username == "" || password == "" || version == "" {
			MissingParameterError(w)
			return
		}

		// Determine the authentication method based on the API version
		if isVersionAbove1_13(version) {
			// Version >= 1.13.0: Token-based authentication
			if err := authenticateWithToken(username, salt, token); err != nil {
				WrongCredentialsError(w)
				return
			}
		} else {
			// Version < 1.13.0: Clear-text or hex-encoded password
			if err := authenticateWithPassword(username, password); err != nil {
				WrongCredentialsError(w)
				return
			}
		}

		// If authentication is successful, call the next handler
		next.ServeHTTP(w, r)
	})
}

func isVersionAbove1_13(version string) bool {
	return strings.Compare(version, "1.13.0") >= 0
}

func authenticateWithToken(username, salt, token string) error {
	user, err := services.GetUserByUsername(username)
	if err != nil || user == nil {
		return errors.New("invalid user")
	}

	// Calculate the expected token using MD5(password + salt)
	expectedToken := calculateMD5Hash(user.Password + salt)

	// Check if the token matches the calculated token
	if token != expectedToken {
		return errors.New("invalid token")
	}

	return nil
}

func authenticateWithPassword(username, password string) error {
	// Check if the password starts with "enc:", meaning it's hex-encoded
	if strings.HasPrefix(password, "enc:") {
		hexPassword := strings.TrimPrefix(password, "enc:")
		passwordBytes, err := hex.DecodeString(hexPassword)
		if err != nil {
			return errors.New("invalid hex-encoded password")
		}
		password = string(passwordBytes)
	}

	user, err := services.GetUserByUsername(username)
	if err != nil || user == nil {
		return errors.New("invalid user")
	}

	// Check if the clear-text password matches
	if user.Password != password {
		return errors.New("invalid password")
	}

	return nil
}

func calculateMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
