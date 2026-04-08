package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductCategory struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	OrganizationID int64          `gorm:"not null;index"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Description    *string        `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"not null;default:now()"`
	UpdatedAt      time.Time      `gorm:"not null;default:now()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (ProductCategory) TableName() string {
	return "product_categories"
}
