package repositories

import (
	"github.com/go-pg/pg/v10"
	"github.com/sakuraapp/api/models"
)

type UserRepository struct {
	db *pg.DB
}

func (r *UserRepository) GetWithDiscriminator(id int64) (*models.User, error) {
	user := new(models.User)
	err := r.db.Model(user).
		Column("user.*").
		ColumnExpr("discriminator.value AS discriminator").
		Join("LEFT JOIN discriminators AS discriminator ON discriminator.owner_id = ?", pg.Ident("user.id")).
		Where("? = ?", pg.Ident("user.id"), id).
		First()

	if err == pg.ErrNoRows {
		err = nil
		user = nil
	}

	return user, err
}

func (r *UserRepository) GetByExternalIdWithDiscriminator(id string) (*models.User, error) {
	user := new(models.User)
	err := r.db.Model(user).
		Column("user.*").
		ColumnExpr("discriminator.value AS discriminator").
		Join("LEFT JOIN discriminators AS discriminator ON discriminator.owner_id = ?", pg.Ident("user.id")).
		Where("? = ?", pg.Ident("user.external_user_id"), id).
		First()

	if err == pg.ErrNoRows {
		err = nil
		user = nil
	}

	return user, err
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.db.Model(user).Insert()

	return err
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := r.db.Model(user).
		WherePK().
		Update()

	return err
}