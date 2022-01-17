package server

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/sakuraapp/api/config"
	"github.com/sakuraapp/api/repository"
	"github.com/sakuraapp/api/store"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const sessionMaxAge = 60 * 60 // 1 hour max - this isn't actually used for sessions, just the sign-up

type Server struct {
	config.Config
	db *pg.DB
	repos *repository.Repositories
	jwt *jwtauth.JWTAuth
	rdb *redis.Client
	cache *cache.Cache
	store store.Service
}

func (s *Server) GetConfig() *config.Config {
	return &s.Config
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

func (s *Server) GetStore() store.Service {
	return s.store
}

func Create(conf config.Config) Server {
	cookieStore := sessions.NewCookieStore([]byte(conf.SessionSecret))
	cookieStore.MaxAge(sessionMaxAge)

	cookieStore.Options.Path = "/"
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.Secure = !conf.IsDev() // only use secure in production

	gothic.Store = cookieStore

	jwtAuth := jwtauth.New("RS256", conf.JWTPrivateKey, conf.JWTPublicKey)

	opts := pg.Options{
		User: conf.DatabaseUser,
		Password: conf.DatabasePassword,
		Database: conf.DatabaseName,
	}

	db := pg.Connect(&opts)
	ctx := context.Background()

	if conf.IsDev() {
		db.AddQueryHook(pgdebug.DebugHook{
			// Print all queries.
			Verbose: true,
		})
	}

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

	myStore := store.NewS3Adapter(&store.S3Config{
		Bucket: conf.S3Bucket,
		Region: conf.S3Region,
		Endpoint: conf.S3Endpoint,
		ForcePathStyle: conf.S3ForcePathStyle,
	})

	repos := repository.Init(db, myCache, myStore)
	s := Server{
		Config: conf,
		db: db,
		repos: &repos,
		jwt: jwtAuth,
		rdb: rdb,
		cache: myCache,
		store: myStore,
	}

	r := NewRouter(&s)

	log.Printf("Listening on port %v", conf.Port)

	err := http.ListenAndServe("0.0.0.0:" + conf.Port, r)

	if err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}

	return s
}