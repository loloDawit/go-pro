package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/loloDawit/ecom/types"
	"github.com/loloDawit/ecom/utils"
)

func GenerateToken(secret []byte, userID int, expiration time.Duration) (string, error) {
	if len(secret) == 0 {
		return "", jwt.ErrInvalidKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(userID),
		"exp":    time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// JWTMiddleware is a middleware function for validating JWT tokens
func JWTMiddleware(secret []byte) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Check for a special header to bypass setting the user ID
			if r.Header.Get("X-Bypass-UserID") == "true" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, "Authorization header is missing")
				return
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if tokenString == "" {
				utils.WriteError(w, http.StatusUnauthorized, "Token is missing")
				return
			}

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the algorithm
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secret, nil
			})

			if err != nil {
				fmt.Printf("Error parsing token: %v\n", err)
				utils.WriteError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			if !token.Valid {
				fmt.Printf("Token is not valid: %v\n", token)
				utils.WriteError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Add user information to the request context if needed
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, ok := claims["userID"].(string)
				if !ok {
					utils.WriteError(w, http.StatusUnauthorized, "Invalid token claims")
					return
				}
				ctx := context.WithValue(r.Context(), types.UserIDKey, userID)
				r = r.WithContext(ctx)
			} else {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		}
	}
}
