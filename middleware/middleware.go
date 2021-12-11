package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sakuraapp/api/internal"
	"github.com/sakuraapp/api/resource"
	apiResource "github.com/sakuraapp/api/resource"
	"github.com/sakuraapp/shared/constant"
	"github.com/sakuraapp/shared/model"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const UserCtxKey = "user"
const SessionCtxKey = "session"

func SendUnauthorized(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, resource.ErrUnauthorized)
}

func SendInternalError(w http.ResponseWriter, r *http.Request)  {
	render.Render(w, r, resource.ErrInternalError)
}

func UserFromContext(ctx context.Context) *model.User {
	user, _ := ctx.Value(UserCtxKey).(*model.User)

	return user
}

func SessionFromContext(ctx context.Context) *internal.Session {
	sess, _ := ctx.Value(SessionCtxKey).(*internal.Session)

	return sess
}

func Authenticator(a internal.App) func(next http.Handler) http.Handler {
	userRepo := a.GetRepositories().User

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx := r.Context()
			token, _, err := jwtauth.FromContext(reqCtx)

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
				SendUnauthorized(w, r)
				return
			}

			user, err := userRepo.GetWithDiscriminator(reqCtx, id)

			if err != nil {
				log.
					WithField("user_id", id).
					WithError(err).
					Error("Failed to get user")

				SendInternalError(w, r)
				return
			}

			if user == nil {
				SendUnauthorized(w, r)
				return
			}

			ctx := context.WithValue(reqCtx, UserCtxKey, user)
			r = r.WithContext(ctx)

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}

func RoomMemberCheck(a internal.App) func(next http.Handler) http.Handler {
	rdb := a.GetRedis()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			strRoomId := chi.URLParam(r, "roomId")
			intRoomId, err := strconv.Atoi(strRoomId)

			if err != nil {
				render.Render(w, r, apiResource.ErrBadRequest)
				return
			}

			roomId := model.RoomId(intRoomId)

			reqCtx := r.Context()

			user := UserFromContext(reqCtx)
			sessionId := r.Header.Get("X-Session-Id")

			if sessionId == "" {
				render.Render(w, r, apiResource.ErrForbidden)
				return
			}

			sessionKey := fmt.Sprintf(constant.SessionFmt, sessionId)

			var sess internal.Session

			err = rdb.HMGet(reqCtx, sessionKey, "user_id", "room_id", "node_id").Scan(&sess)

			if err != nil {
				log.
					WithField("session_id", sessionId).
					WithError(err).
					Error("Failed to retrieve session")

				SendInternalError(w, r)
				return
			}

			if sess.UserId != user.Id || sess.RoomId != roomId {
				render.Render(w, r, apiResource.ErrForbidden)
				return
			}

			sess.Id = sessionId

			ctx := context.WithValue(reqCtx, SessionCtxKey, &sess)
			r = r.WithContext(ctx)

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}