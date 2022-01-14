package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/sakuraapp/shared/model"
)

type RoleRepository struct {
	db *pg.DB
}

func (r *RoleRepository) Get(userId model.UserId, roomId model.RoomId) ([]model.UserRole, error) {
	var roles []model.UserRole
	err := r.db.Model(&roles).
		Column("id", "role_id").
		Where("user_id = ?", userId).
		Where("room_id = ?", roomId).
		Order("id ASC").
		Select()

	if err == pg.ErrNoRows {
		err = nil
	}

	return roles, err
}

func (r *RoleRepository) Add(userRole *model.UserRole) error {
	_, err := r.db.Model(userRole).Insert()

	return err
}