package entity

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          int64  `gorm:"primaryKey;type:bigint;autoIncrement"`
	Name        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Permission) TableName() string {
	return "permissions"
}
