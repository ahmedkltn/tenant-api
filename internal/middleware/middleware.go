package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ahmedkltn/tenant-api/internal/auth"
	"github.com/ahmedkltn/tenant-api/internal/store"
	"github.com/go-chi/chi/v5"
)

// here we have 3 middlewares
// 1. AuthMiddleWare --> is valid JWT ?
// 2. TenantGuard --> JWT.tenant_id == URL param ?
// 3. RoleCheck --> method allowed to role ?

type contextKey string

const claimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TenantGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(claimsKey).(*auth.Claims)
		urlTenantID := chi.URLParam(r, "tenantId")

		if claims.TenantID != urlTenantID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RoleCheck(allowed store.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(claimsKey).(*auth.Claims)

			if claims.Role != allowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
