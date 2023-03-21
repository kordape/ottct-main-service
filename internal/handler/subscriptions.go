package handler

import (
	"fmt"

	"github.com/kordape/ottct-main-service/pkg/logger"
)

type SubscriptionManager struct {
	storage SubscriptionStorage
	log     logger.Interface
}

func NewSubscriptionManager(storage SubscriptionStorage, log logger.Interface) *SubscriptionManager {
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

func (m SubscriptionManager) UpdateSubscription(userId uint, entityId string, subscribe bool) error {
	if subscribe {
		err := m.storage.AddSubscription(userId, entityId)
		if err != nil {
			m.log.Error(fmt.Errorf("[SubscriptionsManager] Failed to update subscription: %w", err))
			return fmt.Errorf("[SubscriptionsManager] Subscription storage error: %w", err)
		}

		return nil
	}

	err := m.storage.DeleteSubscription(userId, entityId)
	if err != nil {
		m.log.Error(fmt.Errorf("[SubscriptionsManager] Failed to update subscription: %w", err))
		return fmt.Errorf("[SubscriptionsManager] Subscription storage error: %w", err)
	}

	return nil
}
