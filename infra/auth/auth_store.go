package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/hjoshi123/WaaS/config"
	"github.com/hjoshi123/WaaS/datastore"
	"github.com/hjoshi123/WaaS/infra/database"
	"github.com/hjoshi123/WaaS/model"
	"github.com/volatiletech/authboss/v3"
)

type AuthStore struct {
	UserStore datastore.UserStore
}

func (a AuthStore) Load(ctx context.Context, key string) (authboss.User, error) {
	// c := authhelper.GinCtx(ctx)
	db := database.Instance(ctx)
	var user model.User

	key = strings.TrimSpace(key)

	//need to check key to see if email or phone
	if strings.Contains(key, "@") {
		if err := db.Where("lower(email) = ?", strings.ToLower(key)).First(&user).Error; err != nil {
			return nil, authboss.ErrUserNotFound
		}
	}
	return &user, nil
}

// Persist the user to DB.
func (a AuthStore) Save(ctx context.Context, user authboss.User) error {
	db := database.Instance(ctx)
	u := user.(*model.User)
	if err := db.Save(u).Error; err != nil {
		return err
	}
	return nil
}

func (a AuthStore) Create(ctx context.Context, user authboss.User) error {
	u, ok := ctx.Value(config.CurrentUser).(*model.User)
	if u != nil && ok {
		return errors.New("cannot register with user already signed in")
	}
	db := database.Instance(ctx)
	u = user.(*model.User)

	//check if email/phone is present
	var check model.User

	//create by email
	if u.Email != "" {
		if err := db.Where("lower(email) = ?", strings.ToLower(u.Email)).First(&check).Error; err == nil {
			return authboss.ErrUserFound
		}
	}

	if err := db.Create(u).Error; err != nil {
		return err
	}
	// Post user create flows.
	// Publish that user was create
	return nil
}

// CreatingServerStorer
func (a AuthStore) New(ctx context.Context) authboss.User {
	return &model.User{}
}

func NewAuthStore() *AuthStore {
	sa := new(AuthStore)
	sa.UserStore = datastore.NewUserStore()

	return sa
}
