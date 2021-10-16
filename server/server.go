package server

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/sakuraapp/api/internal/utils"
	"github.com/sakuraapp/api/repository"
	shared "github.com/sakuraapp/shared/pkg"
	"log"
	"net/http"
	"os"
	"strings"
)

type App struct {
	DB *pg.DB
	Repositories *repository.Repositories
	JWT *jwtauth.JWTAuth
}

func (a *App) GetDB() *pg.DB {
	return a.DB
}

func (a *App) GetRepositories() *repository.Repositories {
	return a.Repositories
}

func (a *App) GetJWT() *jwtauth.JWTAuth {
	return a.JWT
}

func Start(port string) App {
	goth.UseProviders(
		discord.New(
			os.Getenv("DISCORD_KEY"),
			os.Getenv("DISCORD_SECRET"),
			os.Getenv("DISCORD_REDIRECT"),
			GetScopes("DISCORD_SCOPES")...
		),
	)

	// use a fake store because this is a REST API, it's not vulnerable to CSRF anyway
	// todo: re-evaluate this decision
	gothic.Store = utils.NewFakeStore()

	jwtPublicPath := os.Getenv("JWT_PUBLIC_KEY")
	jwtPrivatePath := os.Getenv("JWT_PRIVATE_KEY")
	jwtPassphrase := os.Getenv("JWT_PASSPHRASE")

	jwtPrivateKey := shared.LoadRSAPrivateKey(jwtPrivatePath, jwtPassphrase)
	jwtPublicKey := shared.LoadRSAPublicKey(jwtPublicPath)

	jwtAuth := jwtauth.New("RS256", jwtPrivateKey, jwtPublicKey)

	opts := pg.Options{
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
	}

	db := pg.Connect(&opts)
	ctx := context.Background()

	db.AddQueryHook(pgdebug.DebugHook{
		// Print all queries.
		Verbose: true,
	})

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	repos := repository.Init(db)
	a := App{
		db,
		&repos,
		jwtAuth,
	}
	r := NewRouter(&a)

	log.Printf("Listening on port %v", port)

	err := http.ListenAndServe("0.0.0.0:" + port, r)

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	return a
}

func GetScopes(key string) []string {
	return strings.Split(os.Getenv(key), ", ")
}