package entity

import (
	"time"

	"gorm.io/gorm"
)

type OrganizationCustomer struct {
	ID             int64          `gorm:"primaryKey;type:bigint;autoIncrement"`
	OrganizationID int64          `gorm:"type:bigint;not null"`
	CustomerID     int64          `gorm:"type:bigint;not null"`
	Description    string         `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"not null;default:now()"`
	UpdatedAt      time.Time      `gorm:"not null;default:now()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (OrganizationCustomer) TableName() string {
	return "organization_customers"
}
