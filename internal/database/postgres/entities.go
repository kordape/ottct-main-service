package postgres

import (
	"errors"
	"fmt"

	"github.com/kordape/ottct-main-service/internal/handler"
	"gorm.io/gorm"
)

func (db *DB) GetEntities() ([]handler.Entity, error) {
	var persistentEntities []Entity
	err := db.db.Find(&persistentEntities).Error
	if err != nil {
		return nil, fmt.Errorf("Error getting entities from db: %w", err)
	}

	entities := make([]handler.Entity, len(persistentEntities))
	for i, e := range persistentEntities {
		entities[i] = handler.Entity{
			Id:          e.ID,
			TwitterId:   e.TwitterId,
			DisplayName: e.DisplayName,
		}
	}

	return entities, nil
}

func (db *DB) GetEntity(entityId string) (*handler.Entity, error) {
	var entity Entity

	err := db.db.First(&entity, entityId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("Error getting entities from db: %w", err)
	}

	return &handler.Entity{
		Id:          entity.ID,
		TwitterId:   entity.TwitterId,
		DisplayName: entity.DisplayName,
	}, nil
}
