package repositories

import "github.com/go-pg/pg/v10"

type Repositories struct {
	User UserRepository
	Discriminator DiscriminatorRepository
	Room RoomRepository
}

func Init(db *pg.DB) Repositories {
	return Repositories{
		UserRepository{db},
		DiscriminatorRepository{db},
		RoomRepository{db},
	}
}