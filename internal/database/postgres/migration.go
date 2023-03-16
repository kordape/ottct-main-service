package postgres

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func (db *DB) Migrate() error {
	m := gormigrate.New(db.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202303140000",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&Entity{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("entities")
			},
		},
		{
			ID: "202302240000",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		db.log.Error(fmt.Errorf("Could not migrate: %v", err))
		return fmt.Errorf("Migration failed: %v", err)
	}

	db.log.Info("Migration run successfully")

	return nil
}
