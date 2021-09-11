package internal

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-pg/pg/v10"
	"github.com/sakuraapp/api/repositories"
)

type App interface {
	GetDB() *pg.DB
	GetRepositories() *repositories.Repositories
	GetJWT() *jwtauth.JWTAuth
}