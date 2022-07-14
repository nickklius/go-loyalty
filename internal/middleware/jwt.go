package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nickklius/go-loyalty/config"
	"github.com/nickklius/go-loyalty/internal/utils"
)

type contextType string

const (
	UserIDCtx contextType = "userIDCtx"
)

func JWTMiddleware(cfg config.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, err := utils.ValidateToken(r, &cfg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			userID := token.Claims.(jwt.MapClaims)["user_id"]
			ctx := setClaims(r.Context(), userID.(string))
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func setClaims(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDCtx, userID)
}

func GetClaims(ctx context.Context) string {
	return ctx.Value(UserIDCtx).(string)
}
