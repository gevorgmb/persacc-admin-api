package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64 `gorm:"primaryKey;type:bigint;autoIncrement"`
	Name      string
	Email     string `gorm:"not null;uniqueIndex"`
	RoleID    int64  `gorm:"type:bigint;index"`
	Role      Role   `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
