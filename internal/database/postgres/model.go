package postgres

import (
	"time"
)

const TwitterEntitiesTableName = "twitter_entities"

// User model definition
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Phone     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// TwitterEntity model definition
type TwitterEntity struct {
	ID               string `gorm:"primaryKey"`
	TwitterAccountId string `gorm:"uniqueIndex;not null"`
	DisplayName      string `gorm:"not null"`
}

func (TwitterEntity) TableName() string {
	return TwitterEntitiesTableName
}
