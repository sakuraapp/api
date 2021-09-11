package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/controllers"
	"github.com/sakuraapp/api/internal"
	"github.com/sakuraapp/api/middlewares"
)

func Init(a internal.App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cache-Control"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	/* r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")

			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Cache-Control") // You can add more headers here if needed
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}) */

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	c := controllers.Init(a)

	r.Route("/api/v1", func(r chi.Router) {
		// authentication routes
		r.Route("/auth/{provider}", func(r chi.Router) {
			r.Get("/", c.Auth.BeginAuth)
			r.Get("/callback", c.Auth.CompleteAuth)
		})

		// authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(a.GetJWT()))
			r.Use(middlewares.Authenticator(a))

			// user routes
			r.Get("/users/@me", c.User.GetMyUser)

			// room routes
			r.Route("/rooms", func(r chi.Router) {
				r.Post("/", c.Room.Create)
				r.Get("/{roomId}", c.Room.Get)
				r.Get("/latest", c.Room.GetLatest)
			})
		})
	})

	return r
}