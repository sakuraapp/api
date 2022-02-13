package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/internal/app"
	"github.com/sakuraapp/shared/pkg/constant"
	"github.com/sakuraapp/shared/pkg/model"
	"github.com/sakuraapp/shared/pkg/resource"
	"github.com/sakuraapp/shared/pkg/resource/permission"
	"github.com/sakuraapp/shared/pkg/resource/role"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const UserIdCtxKey = "user_id"
const SessionCtxKey = "session"

func SendUnauthorized(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, resource.ErrUnauthorized)
}

func SendInternalError(w http.ResponseWriter, r *http.Request)  {
	render.Render(w, r, resource.ErrInternalError)
}

func UserIdFromContext(ctx context.Context) model.UserId {
	userId, _ := ctx.Value(UserIdCtxKey).(model.UserId)

	return userId
}

func SessionFromContext(ctx context.Context) *app.Session {
	sess, _ := ctx.Value(SessionCtxKey).(*app.Session)

	return sess
}

func UserValidator(a app.App) func(next http.Handler) http.Handler {
	userRepo := a.GetRepositories().User

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx := r.Context()
			id := UserIdFromContext(reqCtx)

			exists, err := userRepo.Exists(id)

			if err != nil {
				log.
					WithField("user_id", id).
					WithError(err).
					Error("Failed to validate user's existence")

				SendInternalError(w, r)
				return
			}

			if !exists {
				SendUnauthorized(w, r)
				return
			}

			// User exists, pass the request through
			next.ServeHTTP(w, r)
		})
	}
}

func RoomMemberCheck(a app.App) func(next http.Handler) http.Handler {
	rdb := a.GetRedis()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			strRoomId := chi.URLParam(r, "roomId")
			intRoomId, err := strconv.Atoi(strRoomId)

			if err != nil {
				render.Render(w, r, resource.ErrBadRequest)
				return
			}

			roomId := model.RoomId(intRoomId)

			reqCtx := r.Context()

			userId := UserIdFromContext(reqCtx)
			sessionId := r.Header.Get("X-Session-Id")

			if sessionId == "" {
				render.Render(w, r, resource.ErrForbidden)
				return
			}

			sessionKey := fmt.Sprintf(constant.SessionFmt, sessionId)

			var sess app.Session

			err = rdb.HMGet(reqCtx, sessionKey, "user_id", "room_id", "node_id").Scan(&sess)

			if err != nil {
				log.
					WithField("session_id", sessionId).
					WithError(err).
					Error("Failed to retrieve session")

				SendInternalError(w, r)
				return
			}

			if sess.UserId != userId || sess.RoomId != roomId {
				render.Render(w, r, resource.ErrForbidden)
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

func PermissionCheck(perm permission.Permission, a app.App) func(next http.Handler) http.Handler {
	roleRepo := a.GetRepositories().Role

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx := r.Context()
			sess := SessionFromContext(reqCtx)

			userId := sess.UserId
			roomId := sess.RoomId

			userRoles, err := roleRepo.Get(userId, roomId)

			if err != nil {
				log.
					WithError(err).
					WithFields(log.Fields{
						"user_id": userId,
						"room_id": roomId,
					}).
					Error("Failed to retrieve user roles")

				SendInternalError(w, r)
				return
			}

			roles := role.NewManager()

			for _, userRole := range userRoles {
				roles.Add(userRole.RoleId)
			}

			if !roles.HasPermission(perm) {
				render.Render(w, r, resource.ErrForbidden)
				return
			}

			// User has the required permission, so pass their request
			next.ServeHTTP(w, r)
		})
	}
}