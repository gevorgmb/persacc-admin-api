package service

import (
	"context"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) Create(ctx context.Context, product *entity.Product) error {
	return s.DB.Create(product).Error
}

func (s *ProductService) Get(ctx context.Context, id int64, organizationID int64) (*entity.Product, error) {
	var product entity.Product
	err := s.DB.Preload("ProductDetails").Where("id = ? AND organization_id = ?", id, organizationID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) Update(ctx context.Context, product *entity.Product, organizationID int64) error {
	// Verify relationship exists and it matches organization
	var count int64
	s.DB.Model(&entity.Product{}).
		Where("id = ? AND organization_id = ?", product.ID, organizationID).
		Count(&count)
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return s.DB.Save(product).Error
}

func (s *ProductService) Delete(ctx context.Context, id int64, organizationID int64) error {
	return s.DB.Where("id = ? AND organization_id = ?", id, organizationID).Delete(&entity.Product{}).Error
}

func (s *ProductService) List(ctx context.Context, limit, offset int, organizationID int64, filters map[string]string) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := s.DB.Model(&entity.Product{}).Where("organization_id = ?", organizationID)

	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if sku, ok := filters["sku"]; ok && sku != "" {
		query = query.Where("sku ILIKE ?", "%"+sku+"%")
	}
	if description, ok := filters["description"]; ok && description != "" {
		query = query.Where("description ILIKE ?", "%"+description+"%")
	}

	query.Count(&total)
	if err := query.Preload("ProductDetails").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
