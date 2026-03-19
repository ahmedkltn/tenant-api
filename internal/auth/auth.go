package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ahmedkltn/tenant-api/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string     `json:"user_id"`
	TenantID string     `json:"tenant_id"`
	Role     store.Role `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("SfGCo36oyXrFqaICcsdVG6")

func GenerateToken(user *store.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func NewLoginHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, ok := s.FindUserByEmail(req.Email)
		if !ok || user.Password != req.Password {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := GenerateToken(user)
		if err != nil {
			http.Error(w, "could not generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
