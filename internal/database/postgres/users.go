package postgres

import (
	"errors"

	"github.com/kordape/ottct-main-service/internal/handler"
	"gorm.io/gorm"
)

func (db *DB) CreateUser(user handler.User) error {
	u := User{
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
	}

	err := db.db.Create(&u).Error

	return err
}

func (db *DB) GetUserByCredentials(email string, password string) (handler.User, error) {
	u := User{}

	err := db.db.Where("email = ? AND password >= ?", email, password).First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handler.User{}, handler.ErrUserNotFound
	}

	return handler.User{
		Email:    u.Email,
		Password: u.Password,
		Phone:    u.Phone,
	}, err
}
