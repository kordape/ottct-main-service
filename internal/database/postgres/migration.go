package postgres

import (
	_ "embed"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	model "github.com/kordape/ottct-main-service/pkg/db"
)

var (
	//go:embed seed/202303282200.sql
	seed202303282200 string
)

func (db *DB) Migrate(log *logrus.Entry) error {
	m := gormigrate.New(db.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1-entity-schema",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&model.Entity{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("entities")
			},
		},
		{
			ID: "2-user-schema",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&model.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		{
			ID: "3-initial-seed",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(seed202303282200).Error
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Error(fmt.Errorf("Could not migrate: %v", err))
		return fmt.Errorf("Migration failed: %v", err)
	}

	log.Info("Migration run successfully")

	return nil
}
