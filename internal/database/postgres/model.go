package postgres

import (
	"time"
)

// User model definition
type User struct {
	ID            uint      `gorm:"primaryKey"`
	Email         string    `gorm:"uniqueIndex;not null"`
	Password      string    `gorm:"not null"`
	Phone         string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	Subscriptions []Entity  `gorm:"many2many:subscriptions;"`
}

// Entity model definition
type Entity struct {
	ID          string `gorm:"primaryKey"`
	TwitterId   string `gorm:"uniqueIndex;not null"`
	DisplayName string `gorm:"not null"`
}
