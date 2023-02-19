package postgres

import "github.com/kordape/ottct-main-service/internal/handler"

func (db *DB) CreateUser(user handler.User) error {
	u := User{
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
	}

	err := db.db.Create(&u).Error

	return err
}
