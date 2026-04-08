package entity

type ProductDetail struct {
	ID                int64                  `gorm:"primaryKey;autoIncrement"`
	ProductID         int64                  `gorm:"not null;index"`
	AdditionalDetails map[string]interface{} `gorm:"type:jsonb;serializer:json"`
}

func (ProductDetail) TableName() string {
	return "product_details"
}
