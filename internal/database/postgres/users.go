package postgres

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kordape/ottct-main-service/internal/handler"
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

	err := db.db.Where("email = ? AND password = ?", email, password).First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handler.User{}, handler.ErrUserNotFound
	}

	return handler.User{
		Id:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Phone:    u.Phone,
	}, err
}
