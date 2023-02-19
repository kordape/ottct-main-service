package handler

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/logger"
)

var ErrInvalidRequest error = errors.New("Invalid request")

type User struct {
	Email    string
	Password string
	Phone    string
}

type UserStorage interface {
	CreateUser(user User) error
}

type AuthManager struct {
	storage          UserStorage
	log              logger.Interface
	requestValidator *validator.Validate
}

func NewAuthManager(userStorage UserStorage, log logger.Interface, validator *validator.Validate) (AuthManager, error) {

	m := AuthManager{
		storage:          userStorage,
		log:              log,
		requestValidator: validator,
	}

	err := m.validate()

	if err != nil {
		return m, fmt.Errorf("[AuthManager] validation error: %w", err)
	}

	return m, nil
}

func (m AuthManager) validate() error {
	if m.storage == nil {
		return errors.New("user storage is nil")
	}

	return nil
}

func (m AuthManager) SignUp(request api.SignUpRequest) error {
	err := m.validate()

	if err != nil {
		return fmt.Errorf("[AuthManager] manager validation error: %w", err)
	}

	err = m.requestValidator.Struct(request)
	if err != nil {
		m.log.Error(fmt.Errorf("[AuthManager] Invalid SignUp request: %w", err))
		return ErrInvalidRequest
	}

	m.log.Info("[AuthManager] Creating user")

	err = m.storage.CreateUser(User{
		Email:    request.Email,
		Password: request.Password,
		Phone:    request.Phone,
	})

	if err != nil {
		m.log.Error(fmt.Errorf("[AuthManager] Failed to create user: %w", err))
		return fmt.Errorf("[AuthManager] storage error: %w", err)
	}

	m.log.Info("[AuthManager] Created user")

	return nil
}
