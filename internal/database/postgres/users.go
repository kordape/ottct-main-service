package postgres

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kordape/ottct-main-service/internal/handler"
	model "github.com/kordape/ottct-main-service/pkg/db"
)

func (db *DB) CreateUser(user handler.User) (handler.User, error) {
	u := model.User{
		Email:    user.Email,
		Password: user.Password,
	}

	err := db.db.Create(&u).Error

	return handler.User{
		Id:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (db *DB) GetUserByCredentials(email string, password string) (handler.User, error) {
	u := model.User{}

	err := db.db.Where("email = ? AND password = ?", email, password).First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handler.User{}, handler.ErrUserNotFound
	}

	return handler.User{
		Id:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (db *DB) GetUserByEmail(email string) (handler.User, error) {
	u := model.User{}

	err := db.db.Where("email = ?", email).First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handler.User{}, handler.ErrUserNotFound
	}

	return handler.User{
		Id:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}, err
}
