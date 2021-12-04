package internal

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/sakuraapp/api/repository"
	"github.com/sakuraapp/shared/model"
)

type App interface {
	GetDB() *pg.DB
	GetRepositories() *repository.Repositories
	GetJWT() *jwtauth.JWTAuth
	GetRedis() *redis.Client
}

type Session struct {
	Id string `json:"id" redis:"id"`
	UserId model.UserId `json:"user_id" redis:"user_id"`
	RoomId model.RoomId `json:"room_id" redis:"room_id"`
	NodeId string `json:"node_id" redis:"node_id"`
}