package handler

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-main-service/pkg/token"
)

var ErrInvalidRequest error = errors.New("Invalid request")
var ErrUserNotFound error = errors.New("User not found")

const defaultIssuer = "ottct"

type User struct {
	Email    string
	Password string
	Phone    string
}

type UserStorage interface {
	CreateUser(user User) error
	GetUserByCredentials(email string, password string) (User, error)
}

type AuthManager struct {
	storage          UserStorage
	log              logger.Interface
	requestValidator *validator.Validate
	tokenManager     *token.Manager
}

func NewAuthManager(userStorage UserStorage, log logger.Interface, validator *validator.Validate, tokenManager *token.Manager) (AuthManager, error) {

	m := AuthManager{
		storage:          userStorage,
		log:              log,
		requestValidator: validator,
		tokenManager:     tokenManager,
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

	if m.tokenManager == nil {
		return errors.New("token manager is nil")
	}

	if m.requestValidator == nil {
		return errors.New("request validator is nil")
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

	err = m.storage.CreateUser(User{
		Email:    request.Email,
		Password: request.Password,
		Phone:    request.Phone,
	})

	if err != nil {
		m.log.Error(fmt.Errorf("[AuthManager] Failed to create user: %w", err))
		return fmt.Errorf("[AuthManager] storage error: %w", err)
	}

	return nil
}

func (m AuthManager) Auth(request api.AuthRequest) (string, error) {
	err := m.validate()

	if err != nil {
		return "", fmt.Errorf("[AuthManager] manager validation error: %w", err)
	}

	err = m.requestValidator.Struct(request)
	if err != nil {
		m.log.Error(fmt.Errorf("[AuthManager] Invalid Auth request: %w", err))
		return "", ErrInvalidRequest
	}

	user, err := m.storage.GetUserByCredentials(request.Email, request.Password)

	if err != nil {
		m.log.Error(fmt.Errorf("[AuthManager] Failed to get user: %w", err))
		return "", fmt.Errorf("[AuthManager] storage error: %w", err)
	}

	token, err := m.tokenManager.GenerateJWT(user.Email)
	if err != nil {
		m.log.Error(fmt.Errorf("[AuthManager] Failed to generate token for user: %w", err))
		return "", fmt.Errorf("[AuthManager] token manager error: %w", err)
	}

	return token, nil
}
