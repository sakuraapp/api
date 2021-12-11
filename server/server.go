package server

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/markbates/goth/gothic"
	"github.com/sakuraapp/api/config"
	"github.com/sakuraapp/api/internal/utils"
	"github.com/sakuraapp/api/repository"
	"log"
	"net/http"
)

type Server struct {
	db *pg.DB
	repos *repository.Repositories
	jwt *jwtauth.JWTAuth
	rdb *redis.Client
	cache *cache.Cache
}

func (s *Server) GetDB() *pg.DB {
	return s.db
}

func (s *Server) GetRepositories() *repository.Repositories {
	return s.repos
}

func (s *Server) GetJWT() *jwtauth.JWTAuth {
	return s.jwt
}

func (s *Server) GetRedis() *redis.Client {
	return s.rdb
}

func (s *Server) GetCache() *cache.Cache {
	return s.cache
}

func Create(conf config.Config) Server {
	// use a fake store because this is a REST API, it's not vulnerable to CSRF anyway
	// todo: re-evaluate this decision
	gothic.Store = utils.NewFakeStore()

	jwtAuth := jwtauth.New("RS256", conf.JWTPrivateKey, conf.JWTPublicKey)

	opts := pg.Options{
		User: conf.DatabaseUser,
		Password: conf.DatabasePassword,
		Database: conf.DatabaseName,
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

	rdb := redis.NewClient(&redis.Options{
		Addr: conf.RedisAddr,
		Password: conf.RedisPassword,
		DB: conf.RedisDatabase,
	})

	myCache := cache.New(&cache.Options{
		Redis: rdb,
		// LocalCache: cache.NewTinyLFU(1000, time.Minute),
		// until server-assisted client cache is possible, don't keep a client cache (we can't invalidate it)
	})

	repos := repository.Init(db, myCache)
	s := Server{
		db: db,
		repos: &repos,
		jwt: jwtAuth,
		rdb: rdb,
		cache: myCache,
	}

	r := NewRouter(&s)

	log.Printf("Listening on port %v", conf.Port)

	err := http.ListenAndServe("0.0.0.0:" + conf.Port, r)

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	return s
}