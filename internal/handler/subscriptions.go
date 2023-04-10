package handler

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/sirupsen/logrus"
)

type SubscriptionManager struct {
	storage          SubscriptionStorage
	requestValidator *validator.Validate
}

func NewSubscriptionManager(storage SubscriptionStorage, validate *validator.Validate) (*SubscriptionManager, error) {
	m := &SubscriptionManager{
		storage:          storage,
		requestValidator: validate,
	}

	err := m.validate()
	if err != nil {
		return nil, fmt.Errorf("error validation subscription manager: %w", err)
	}

	return m, nil
}

func (m SubscriptionManager) validate() error {
	if m.storage == nil {
		return errors.New("subscription storage is nil")
	}

	if m.requestValidator == nil {
		return errors.New("request validator is nil")
	}

	return nil
}

type SubscriptionStorage interface {
	GetSubscriptionsByUser(userId uint) ([]Entity, error)
	GetSubscriptionsByEntity(entityId string) ([]User, error)
	AddSubscription(userId uint, entityId string) error
	DeleteSubscription(userId uint, entityId string) error
}

func (m SubscriptionManager) GetSubscriptionsByUser(userId uint, log *logrus.Entry) (entities []Entity, err error) {
	entities, err = m.storage.GetSubscriptionsByUser(userId)
	if err != nil {
		log.WithError(err).Error("[SubscriptionsManager] Failed to get user's subscriptions")
		return nil, fmt.Errorf("[SubscriptionsManager] Subscription storage error: %w", err)
	}

	return
}

func (m SubscriptionManager) GetSubscriptionsByEntity(entityId string, log *logrus.Entry) (users []User, err error) {
	users, err = m.storage.GetSubscriptionsByEntity(entityId)
	if err != nil {
		log.WithError(err).Error("[SubscriptionsManager] Failed to get subscriptions by entity")
		return nil, fmt.Errorf("[SubscriptionsManager] Subscription storage error: %w", err)
	}

	return
}

func (m SubscriptionManager) UpdateSubscription(userId uint, entityId string, request api.UpdateSubscriptionRequest, log *logrus.Entry) error {
	err := m.requestValidator.Struct(request)
	if err != nil {
		log.WithError(err).Error("[SubscriptionManager] Error validating request")
		return ErrInvalidRequest
	}

	if request.Subscribe {
		err := m.storage.AddSubscription(userId, entityId)
		if err != nil {
			log.WithError(err).Error("[SubscriptionsManager] Failed to update subscription")
			return fmt.Errorf("subscription storage error: %w", err)
		}

		return nil
	}

	err = m.storage.DeleteSubscription(userId, entityId)
	if err != nil {
		log.WithError(err).Error("[SubscriptionsManager] Failed to update subscription")
		return fmt.Errorf("subscription storage error: %w", err)
	}

	return nil
}
