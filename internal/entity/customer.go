package entity

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID             int64                  `gorm:"primaryKey;autoIncrement"`
	Name           string                 `gorm:"type:varchar(255);not null"`
	FirstName      string                 `gorm:"type:varchar(255)"`
	LastName       string                 `gorm:"type:varchar(255)"`
	Prefix         string                 `gorm:"type:varchar(255)"`
	MiddleName     string                 `gorm:"type:varchar(255)"`
	Suffix         string                 `gorm:"type:varchar(255)"`
	Birthday       *time.Time             `gorm:"type:date"`
	Phone          string                 `gorm:"type:varchar(255)"`
	Email          string                 `gorm:"type:varchar(255)"`
	AdditionalInfo map[string]interface{} `gorm:"type:jsonb;serializer:json"`
	UserID         *int64                 `gorm:"type:bigint;unique;default:null"`
	CreatedAt      time.Time              `gorm:"not null;default:now()"`
	UpdatedAt      time.Time              `gorm:"not null;default:now()"`
	DeletedAt      gorm.DeletedAt         `gorm:"index"`
}

func (Customer) TableName() string {
	return "customers"
}
