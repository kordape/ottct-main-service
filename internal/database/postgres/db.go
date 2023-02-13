package postgres

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) *gorm.DB {

	db.AutoMigrate(&User{})

	return db
}
