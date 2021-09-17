package auth

import (
	"context"
	"net/http"

	"github.com/gpiecyk/data-warehouse/internal/api"
	"github.com/gpiecyk/data-warehouse/internal/users"
)

var userContextKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware(api *api.API) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			id, err := validateJwtToken(header)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user, err := api.GetUserById(r.Context(), id)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := NewContextWithUser(r.Context(), user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func validateJwtToken(header string) (id int, err error) {
	tokenString := header
	id, err = ParseToken(tokenString)
	return
}

func NewContextWithUser(ctx context.Context, user *users.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) *users.User {
	user, _ := ctx.Value(userContextKey).(*users.User)
	return user
}
