package app

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/sakuraapp/api/internal/adapter"
	"github.com/sakuraapp/api/internal/config"
	"github.com/sakuraapp/api/internal/repository"
	"github.com/sakuraapp/api/pkg/store"
	"github.com/sakuraapp/pubsub"
	"github.com/sakuraapp/shared/pkg/model"
	"github.com/sakuraapp/shared/pkg/resource"
)

type App interface {
	pubsub.Dispatcher
	GetConfig() *config.Config
	GetDB() *pg.DB
	GetRepositories() *repository.Repositories
	GetAdapters() *adapter.Adapters
	GetJWT() *jwtauth.JWTAuth
	GetRedis() *redis.Client
	GetCache() *cache.Cache
	GetStore() store.Service
	GetBuilder() *resource.Builder
}

type Session struct {
	Id     string       `json:"id" redis:"id"`
	UserId model.UserId `json:"user_id" redis:"user_id"`
	RoomId model.RoomId `json:"room_id" redis:"room_id"`
	NodeId string       `json:"node_id" redis:"node_id"`
}
