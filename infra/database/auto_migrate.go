package database

import (
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/model"
	"github.com/waas-app/WaaS/util"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.Device{})
	db.AutoMigrate(&model.User{})

	// Insert an user to the database
	u := new(model.User)
	u.Username = config.Spec.AdminUserName
	u.Email = config.Spec.AdminUserName

	hashedPassword, err := util.HashPassword(config.Spec.AdminPassword)
	if err != nil {
		panic(err)
	}

	u.PutPassword(hashedPassword)
	u.Admin = true

	err = db.Save(u).Error
	if err != nil {
		panic(err)
	}
}
