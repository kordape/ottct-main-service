package postgres

import (
	"fmt"

	"github.com/kordape/ottct-main-service/internal/handler"
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
