package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/logger"
)

type SubscriptionManager struct {
	storage          SubscriptionStorage
	log              logger.Interface
	requestValidator *validator.Validate
}

func NewSubscriptionManager(storage SubscriptionStorage, log logger.Interface, validate *validator.Validate) *SubscriptionManager {
	return &SubscriptionManager{
		storage: storage,
		log:     log,
	}
}

type SubscriptionStorage interface {
	GetSubscriptionsByUser(userId uint) ([]Entity, error)
	AddSubscription(userId uint, entityId string) error
	DeleteSubscription(userId uint, entityId string) error
}

func (m SubscriptionManager) GetSubscriptionsByUser(userId uint) (entities []Entity, err error) {
	entities, err = m.storage.GetSubscriptionsByUser(userId)
	if err != nil {
		m.log.Error(fmt.Errorf("[SubscriptionsManager] Failed to get user's subscriptions: %w", err))
		return nil, fmt.Errorf("[SubscriptionsManager] Subscription storage error: %w", err)
	}

	return
}

func (m SubscriptionManager) UpdateSubscription(userId uint, entityId string, request api.UpdateSubscriptionRequest) error {
	err := m.requestValidator.Struct(request)
	if err != nil {
		m.log.Error(fmt.Errorf("[SubscriptionsManager] Error validating request: %w", err))
		return fmt.Errorf("[SubscriptionsManager]  Error validating request: %w", err)
	}

	if request.Subscribe {
		err := m.storage.AddSubscription(userId, entityId)
		if err != nil {
			m.log.Error(fmt.Errorf("[SubscriptionsManager] Failed to update subscription: %w", err))
			return fmt.Errorf("[SubscriptionsManager] Subscription storage error: %w", err)
		}

		return nil
	}

	err = m.storage.DeleteSubscription(userId, entityId)
	if err != nil {
		m.log.Error(fmt.Errorf("[SubscriptionsManager] Failed to update subscription: %w", err))
		return fmt.Errorf("[SubscriptionsManager] Subscription storage error: %w", err)
	}

	return nil
}
