package datastore

import (
	"context"

	"github.com/waas-app/WaaS/infra/database"
	"github.com/waas-app/WaaS/model"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
)

type UserStore interface {
	SaveUser(ctx context.Context, user *model.User) error
	FindByQuery(ctx context.Context, query string, values ...interface{}) ([]*model.User, error)
	FindUserByID(ctx context.Context, id uint) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

func NewUserStore() UserStore {
	store := new(userStore)

	return store
}

type userStore struct{}

func (userStore) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	db := database.Instance(ctx)
	user := new(model.User)
	if err := db.First(user, "email = ?", email).Error; err != nil {
		util.Logger(ctx).Error("could not find user by email", zap.Error(err), zap.String("email", email))
		return nil, err
	}
	return user, nil
}

func (userStore) FindUserByID(ctx context.Context, id uint) (*model.User, error) {
	db := database.Instance(ctx)
	user := new(model.User)
	if err := db.First(user, "id = ?", id).Error; err != nil {
		util.Logger(ctx).Error("could not find user by id", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}
	return user, nil
}

func (userStore) FindByQuery(ctx context.Context, query string, values ...interface{}) ([]*model.User, error) {
	var users []*model.User
	db := database.Instance(ctx)
	if err := db.Where(query, values...).Find(&users).Error; err != nil {
		util.Logger(ctx).Error("could not find users by query", zap.Error(err), zap.String("query", query))
		return nil, err
	}
	return users, nil
}

func (userStore) SaveUser(ctx context.Context, user *model.User) error {
	db := database.Instance(ctx)

	if err := db.Save(user).Error; err != nil {
		util.Logger(ctx).Error("could not save user", zap.Error(err), zap.String("email", user.Email))
		return err
	}
	return nil
}
