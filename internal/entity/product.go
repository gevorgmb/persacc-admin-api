package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	OrganizationID int64          `gorm:"not null;index"`
	SKU            string         `gorm:"type:varchar(255);not null;uniqueIndex"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Description    *string        `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"not null;default:now()"`
	UpdatedAt      time.Time      `gorm:"not null;default:now()"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	ProductDetails    *ProductDetail `gorm:"foreignKey:ProductID"`
	CategoryID        *int64         `gorm:"index;default:null"`
	VendorID          *int64         `gorm:"index;default:null"`
	VendorProductCode *string        `gorm:"type:varchar(255);default:null"`
}

func (Product) TableName() string {
	return "products"
}
