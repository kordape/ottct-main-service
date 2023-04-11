package handler

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/token"
)

var ErrInvalidRequest error = errors.New("Invalid request")
var ErrUserNotFound error = errors.New("User not found")

const defaultIssuer = "ottct"

type User struct {
	Id       uint
	Email    string
	Password string
}

type UserStorage interface {
	CreateUser(user User) (User, error)
	GetUserByCredentials(email string, password string) (User, error)
	GetUserByEmail(email string) (User, error)
}

type AuthManager struct {
	storage          UserStorage
	requestValidator *validator.Validate
	tokenManager     *token.Manager
}

func NewAuthManager(userStorage UserStorage, validator *validator.Validate, tokenManager *token.Manager) (*AuthManager, error) {

	m := AuthManager{
		storage:          userStorage,
		requestValidator: validator,
		tokenManager:     tokenManager,
	}

	err := m.validate()

	if err != nil {
		return &m, fmt.Errorf("[AuthManager] validation error: %w", err)
	}

	return &m, nil
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

func (m AuthManager) SignUp(request api.SignUpRequest, log *logrus.Entry) (User, error) {
	err := m.validate()

	if err != nil {
		return User{}, fmt.Errorf("[AuthManager] manager validation error: %w", err)
	}

	err = m.requestValidator.Struct(request)
	if err != nil {
		log.WithError(err).Error("[AuthManager] Invalid SignUp request")
		return User{}, ErrInvalidRequest
	}

	user, err := m.storage.CreateUser(User{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		log.WithError(err).Error("[AuthManager] Failed to create user")
		return User{}, fmt.Errorf("[AuthManager] storage error: %w", err)
	}

	return user, nil
}

func (m AuthManager) Auth(request api.AuthRequest, log *logrus.Entry) (string, error) {
	err := m.validate()

	if err != nil {
		return "", fmt.Errorf("[AuthManager] manager validation error: %w", err)
	}

	err = m.requestValidator.Struct(request)
	if err != nil {
		log.WithError(err).Error("[AuthManager] Invalid Auth request")
		return "", ErrInvalidRequest
	}

	user, err := m.storage.GetUserByCredentials(request.Email, request.Password)

	if err != nil {
		log.WithError(err).Error("[AuthManager] Failed to get user")
		return "", fmt.Errorf("[AuthManager] storage error: %w", err)
	}

	token, err := m.tokenManager.GenerateJWT(user.Id)
	if err != nil {
		log.WithError(err).Error("[AuthManager] Failed to generate token for user")
		return "", fmt.Errorf("[AuthManager] token manager error: %w", err)
	}

	return token, nil
}

func (m AuthManager) GetUserByEmail(email string, log *logrus.Entry) (User, error) {
	err := m.validate()

	if err != nil {
		return User{}, fmt.Errorf("[AuthManager] manager validation error: %w", err)
	}

	user, err := m.storage.GetUserByEmail(email)

	if err != nil {
		log.WithError(err).Warn("[AuthManager] Failed to get user")
		return User{}, fmt.Errorf("[AuthManager] storage error: %w", err)
	}

	return user, nil
}
