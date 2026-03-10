package entity

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID          int64          `gorm:"primaryKey;type:bigint;autoIncrement"`
	OwnerID     int64          `gorm:"type:bigint;not null;index"`
	Name        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Owner       User           `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (Organization) TableName() string {
	return "organizations"
}
