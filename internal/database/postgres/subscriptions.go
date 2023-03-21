package postgres

import (
	"fmt"

	"github.com/kordape/ottct-main-service/internal/handler"
)

func (db *DB) GetSubscriptionsByUser(userId uint) ([]handler.Entity, error) {
	var subscriptions []Entity
	err := db.db.Model(&User{ID: userId}).Association("Subscriptions").Find(&subscriptions)
	if err != nil {
		return nil, fmt.Errorf("Error getting user's subscriptions from db: %w", err)
	}

	entities := make([]handler.Entity, len(subscriptions))
	for i, e := range subscriptions {
		entities[i] = handler.Entity{
			Id:          e.ID,
			TwitterId:   e.TwitterId,
			DisplayName: e.DisplayName,
		}
	}

	return entities, nil
}

func (db *DB) AddSubscription(userId uint, entityId string) error {
	err := db.db.Model(&User{ID: userId}).Association("Subscriptions").Append(&Entity{ID: entityId})
	if err != nil {
		return fmt.Errorf("Error adding subscription to the user: %w", err)
	}

	return nil
}

func (db *DB) DeleteSubscription(userId uint, entityId string) error {
	err := db.db.Model(&User{ID: userId}).Association("Subscriptions").Delete(&Entity{ID: entityId})
	if err != nil {
		return fmt.Errorf("Error deleting user's subscription: %w", err)
	}

	return nil
}
