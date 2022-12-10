package database

import (
	"github.com/hjoshi123/WaaS/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.Device{})
}
