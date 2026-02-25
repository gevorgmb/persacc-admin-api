package entity

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          int64        `gorm:"primaryKey;type:bigint;autoIncrement"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Role) TableName() string {
	return "roles"
}
