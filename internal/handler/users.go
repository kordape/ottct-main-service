package handler

import (
	"errors"
	"fmt"

	"github.com/kordape/ottct-main-service/pkg/logger"
)

type User struct {
	Email    string
	Password string
	Phone    string
}

type UserStorage interface {
	CreateUser(user User) error
}

type UserManager struct {
	storage UserStorage
	log     logger.Interface
}

func NewUserManager(userStorage UserStorage, log logger.Interface) (UserManager, error) {

	m := UserManager{
		storage: userStorage,
		log:     log,
	}

	err := m.validate()

	if err != nil {
		return m, fmt.Errorf("[UserManager] validation error: %w", err)
	}

	return m, nil
}

func (m UserManager) validate() error {
	if m.storage == nil {
		return errors.New("user storage is nil")
	}

	return nil
}

func (m UserManager) CreateUser() error {
	err := m.validate()

	if err != nil {
		return fmt.Errorf("[UserManager] manager validation error: %w", err)
	}

	m.log.Info("[UserManager] Creating user")

	err = m.storage.CreateUser(User{
		Email:    "petar.korda@gmail.com",
		Password: "1234",
		Phone:    "0643025572",
	})

	if err != nil {
		m.log.Error("[UserManager] Failed to create user: %w", err)
		return fmt.Errorf("[UserManager] storage error: %w", err)
	}

	return nil
}
