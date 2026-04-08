package entity

import (
	"time"

	"gorm.io/gorm"
)

type Supplier struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Domain         *string        `gorm:"type:varchar(255)"`
	Phone          *string        `gorm:"type:varchar(255)"`
	Description    *string        `gorm:"type:text"`
	OrganizationID *int64         `gorm:"index"`
	CreatedAt      time.Time      `gorm:"not null;default:now()"`
	UpdatedAt      time.Time      `gorm:"not null;default:now()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
