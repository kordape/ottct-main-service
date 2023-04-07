package postgres

import (
	"fmt"

	"github.com/kordape/ottct-main-service/internal/handler"
	model "github.com/kordape/ottct-main-service/pkg/db"
)

func (db *DB) GetSubscriptionsByUser(userId uint) ([]handler.Entity, error) {
	var subscriptions []model.Entity
	err := db.db.Model(&model.User{ID: userId}).Association("Subscriptions").Find(&subscriptions)
	if err != nil {
		return nil, fmt.Errorf("error getting user's subscriptions from db: %w", err)
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
	err := db.db.Model(&model.User{ID: userId}).Association("Subscriptions").Append(&model.Entity{ID: entityId})
	if err != nil {
		return fmt.Errorf("error adding subscription to the user: %w", err)
	}

	return nil
}

func (db *DB) DeleteSubscription(userId uint, entityId string) error {
	err := db.db.Model(&model.User{ID: userId}).Association("Subscriptions").Delete(&model.Entity{ID: entityId})
	if err != nil {
		return fmt.Errorf("error deleting user's subscription: %w", err)
	}

	return nil
}

func (db *DB) GetSubscriptionsByEntity(entityId string) ([]handler.User, error) {
	var subscriptions []model.User
	err := db.db.Model(&model.Entity{ID: entityId}).Association("Subscriptions").Find(&subscriptions)
	if err != nil {
		return nil, fmt.Errorf("error getting subscribed users from db: %w", err)
	}

	users := make([]handler.User, len(subscriptions))
	for i, e := range subscriptions {
		users[i] = handler.User{
			Id:    e.ID,
			Email: e.Email,
		}
	}

	return users, nil
}
