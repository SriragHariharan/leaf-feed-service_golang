package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sriraghariharan/feed-service-go/internal/httputil"
)

const expectedIssuer = "leaf.com"

type accessTokenClaims struct {
	jwt.RegisteredClaims
}

func ValidateAccessToken(next http.Handler) http.Handler {
	secret := os.Getenv("ACCESS_TOKEN_SECRET")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "authorization header is required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "token is required")
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims := &accessTokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid or expired token")
			return
		}

		if claims.Issuer != "" && claims.Issuer != expectedIssuer {
			httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid token issuer")
			return
		}

		userID, ok := audienceUserID(claims.Audience)
		if !ok {
			httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "token is missing user identity")
			return
		}

		next.ServeHTTP(w, r.WithContext(WithUserID(r.Context(), userID)))
	})
}

func audienceUserID(audience jwt.ClaimStrings) (string, bool) {
	if len(audience) == 0 {
		return "", false
	}
	userID := strings.TrimSpace(audience[0])
	if userID == "" {
		return "", false
	}
	return userID, true
}
