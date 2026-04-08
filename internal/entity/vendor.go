package entity

import (
	"time"

	"gorm.io/gorm"
)

type Vendor struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Domain      *string        `gorm:"type:varchar(255)"`
	Description *string        `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Vendor) TableName() string {
	return "vendors"
}
