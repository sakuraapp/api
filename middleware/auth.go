package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sakuraapp/api/internal"
	"github.com/sakuraapp/api/response"
	"github.com/sakuraapp/shared/model"
	"net/http"
)

const UserCtxKey = "user"

func SendUnauthorized(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, response.ErrUnauthorized)
	return
}

func FromContext(ctx context.Context) *model.User {
	user, _ := ctx.Value(UserCtxKey).(*model.User)

	return user
}

func Authenticator(a internal.App) func(next http.Handler) http.Handler {
	userRepo := a.GetRepositories().User

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				SendUnauthorized(w, r)
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				SendUnauthorized(w, r)
				return
			}

			rawId, _ := token.Get("id")
			floatId, _ := rawId.(float64)
			id := model.UserId(floatId)

			if id == 0 {
				fmt.Printf("not valid %v\n", id)
				SendUnauthorized(w, r)
				return
			}

			user, err := userRepo.GetWithDiscriminator(id)

			if user == nil || err != nil {
				if err != nil {
					fmt.Printf("error %v", err.Error())
				}

				SendUnauthorized(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserCtxKey, user)
			r = r.WithContext(ctx)

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}