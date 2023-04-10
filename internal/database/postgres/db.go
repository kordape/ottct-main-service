package postgres

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type DB struct {
	db  *gorm.DB
	log *logrus.Entry
}

func New(db *gorm.DB, log *logrus.Entry) (*DB, error) {
	d := &DB{
		db:  db,
		log: log,
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

	if db.log == nil {
		return errors.New("validation error: log is nil")
	}

	return nil
}
