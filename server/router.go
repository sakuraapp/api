package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/controller"
	"github.com/sakuraapp/api/internal"
	sakuraMiddleware "github.com/sakuraapp/api/middleware"
	sharedMiddleware "github.com/sakuraapp/shared/pkg/middleware"
	"github.com/sakuraapp/shared/pkg/resource/permission"
)

func NewRouter(a internal.App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   a.GetConfig().AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cache-Control", "X-Session-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	c := controller.Init(a)

	r.Route("/v1", func(r chi.Router) {
		// authentication routes
		r.Route("/auth", func(r chi.Router) {
			r.Route("/{provider}", func(r chi.Router) {
				r.Get("/", c.Auth.BeginAuth)
				r.Get("/callback", c.Auth.HandleCallback)
			})

			r.Post("/complete", c.Auth.CompleteAuth)
		})

		// authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(a.GetJWT()))
			r.Use(sharedMiddleware.Authenticator)
			r.Use(sakuraMiddleware.UserValidator(a))

			// user routes
			r.Get("/users/@me", c.User.GetMyUser)

			// room routes
			r.Route("/rooms", func(r chi.Router) {
				r.Post("/", c.Room.Create)
				r.Get("/latest", c.Room.GetLatest)

				r.Route("/{roomId}", func(r chi.Router) {
					r.Get("/", c.Room.Get)
					r.Put("/", c.Room.Update)

					// room routes for room members only
					r.Group(func(r chi.Router) {
						r.Use(sakuraMiddleware.RoomMemberCheck(a))

						r.Get("/queue", c.Room.GetQueue)
						r.Post("/messages", c.Room.SendMessage)

						r.Group(func(r chi.Router) {
							r.Use(sakuraMiddleware.PermissionCheck(permission.MANAGE_ROOM, a))

							r.Post("/vm", c.Room.DeployVM)
						})
					})
				})
			})
		})
	})

	return r
}