package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo_list/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	UserIDKey   = "user_id"
	TenantIDKey = "tenant_id"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string, database *db.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := ""

			// Intentar obtener token de cookie primero
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				tokenString = cookie.Value
			}

			// Si no hay cookie, intentar desde header Authorization
			if tokenString == "" {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "" {
					parts := strings.Split(authHeader, " ")
					if len(parts) == 2 && parts[0] == "Bearer" {
						tokenString = parts[1]
					}
				}
			}

			// Si no hay token, redirigir a login
			if tokenString == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, TenantIDKey, claims.TenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}

func GetTenantID(ctx context.Context) uuid.UUID {
	tenantID, ok := ctx.Value(TenantIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return tenantID
}
