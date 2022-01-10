package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/sakuraapp/shared/model"
)

type RoleRepository struct {
	db *pg.DB
}

func (r *RoleRepository) Add(userRole *model.UserRole) error {
	_, err := r.db.Model(userRole).Insert()

	return err
}