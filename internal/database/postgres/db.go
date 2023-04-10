package postgres

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*DB, error) {
	d := &DB{
		db: db,
	}

	err := d.validate()

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize DB: %v", err)
	}

	return d, nil
}

func (db *DB) validate() error {
	if db.db == nil {
		return errors.New("validation error: db is nil")
	}

	return nil
}
