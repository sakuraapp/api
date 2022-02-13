package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/cache/v8"
	"github.com/sakuraapp/api/pkg/store"
)

type Repositories struct {
	User          UserRepository
	Discriminator DiscriminatorRepository
	Room          RoomRepository
	Role          RoleRepository
}

func Init(db *pg.DB, cache *cache.Cache, store store.Service) Repositories {
	return Repositories{
		User:          UserRepository{db: db, cache: cache, store: store},
		Discriminator: DiscriminatorRepository{db: db},
		Room:          RoomRepository{db: db, cache: cache},
		Role:          RoleRepository{db: db},
	}
}