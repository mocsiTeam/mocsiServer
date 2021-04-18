package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/mocsiTeam/mocsiServer/auth/jwt"
	"github.com/mocsiTeam/mocsiServer/db"
	"gorm.io/gorm"
)

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

var DB = db.Connector()

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("token")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := header
			id, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}
			userID, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}
			user := db.Users{Model: gorm.Model{ID: uint(userID)}}
			if user.Check(DB) != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), UserCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *db.Users {
	raw, _ := ctx.Value(UserCtxKey).(*db.Users)
	return raw
}
