package auth

import (
	"context"
	"strings"

	"github.com/hjoshi123/WaaS/datastore"
	"github.com/hjoshi123/WaaS/infra/database"
	"github.com/hjoshi123/WaaS/model"
	"github.com/volatiletech/authboss"
)

type AuthStore struct {
	UserStore datastore.UserStore
}

func (a AuthStore) Load(ctx context.Context, key string) (authboss.User, error) {
	// c := authhelper.GinCtx(ctx)
	db := database.Instance(ctx)
	var user model.User

	key = strings.TrimSpace(key)

	db = db.Where("active = true or active is null")

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

// CreatingServerStorer
func (a AuthStore) New(ctx context.Context) authboss.User {
	return &model.User{}
}

func NewAuthStore() *AuthStore {
	sa := new(AuthStore)
	sa.UserStore = datastore.NewUserStore()

	return sa
}
