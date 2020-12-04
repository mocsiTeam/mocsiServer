package auth

import (
	"context"
	"net/http"

	"github.com/mocsiTeam/mocsiServer/auth/jwt"
	"github.com/mocsiTeam/mocsiServer/db"
)

var userCtxKey = &contextKey{"user"}

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
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user := db.Users{NickName: username}
			err = user.Check(DB)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *db.Users {
	raw, _ := ctx.Value(userCtxKey).(*db.Users)
	return raw
}
