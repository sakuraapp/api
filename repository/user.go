package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/cache/v8"
	"github.com/sakuraapp/shared/constant"
	"github.com/sakuraapp/shared/model"
)

type UserRepository struct {
	db *pg.DB
	cache *cache.Cache
}

func (u *UserRepository) GetWithDiscriminator(ctx context.Context, id model.UserId) (*model.User, error) {
	user := new(model.User)

	if err := u.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf(constant.UserCacheFmt, id),
		Value: user,
		TTL:   constant.UserCacheTTL,
		Do: func(item *cache.Item) (interface{}, error) {
			return u.fetchWithDiscriminator(user, id)
		},
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) fetchWithDiscriminator(user *model.User, id model.UserId) (*model.User, error) {
	err := u.db.Model(user).
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

func (u *UserRepository) FetchWithDiscriminator(id model.UserId) (*model.User, error) {
	user := new(model.User)

	return u.fetchWithDiscriminator(user, id)
}

func (u *UserRepository) GetByExternalIdWithDiscriminator(id string) (*model.User, error) {
	user := new(model.User)
	err := u.db.Model(user).
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

func (u *UserRepository) GetUsername(id model.UserId) (string, error) {
	var username string
	_, err := u.db.QueryOne(pg.Scan(&username), "SELECT username FROM users WHERE id = ? LIMIT 1", id)

	return username, err
}

func (u *UserRepository) Create(user *model.User) error {
	_, err := u.db.Model(user).Insert()

	return err
}

func (u *UserRepository) Update(user *model.User) error {
	_, err := u.db.Model(user).
		WherePK().
		Update()

	return err
}

func (u *UserRepository) Exists(id model.UserId) (bool, error) {
	var exists bool
	err := u.db.Model(pg.Scan(&exists)).
		ColumnExpr("EXISTS(SELECT 1 FROM users WHERE id = ?)", id).
		Select(&exists)

	return exists, err
}