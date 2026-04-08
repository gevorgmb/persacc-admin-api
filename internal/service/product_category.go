package service

import (
	"context"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type ProductCategoryService struct {
	DB *gorm.DB
}

func NewProductCategoryService(db *gorm.DB) *ProductCategoryService {
	return &ProductCategoryService{DB: db}
}

func (s *ProductCategoryService) Create(ctx context.Context, category *entity.ProductCategory) error {
	return s.DB.Create(category).Error
}

func (s *ProductCategoryService) Get(ctx context.Context, id int64, organizationID int64) (*entity.ProductCategory, error) {
	var category entity.ProductCategory
	err := s.DB.Where("id = ? AND organization_id = ?", id, organizationID).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *ProductCategoryService) Update(ctx context.Context, category *entity.ProductCategory, organizationID int64) error {
	// Verify relationship exists and it matches organization
	var count int64
	s.DB.Model(&entity.ProductCategory{}).
		Where("id = ? AND organization_id = ?", category.ID, organizationID).
		Count(&count)
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return s.DB.Save(category).Error
}

func (s *ProductCategoryService) Delete(ctx context.Context, id int64, organizationID int64) error {
	return s.DB.Where("id = ? AND organization_id = ?", id, organizationID).Delete(&entity.ProductCategory{}).Error
}

func (s *ProductCategoryService) List(ctx context.Context, limit, offset int, organizationID int64, filters map[string]string) ([]entity.ProductCategory, int64, error) {
	var categories []entity.ProductCategory
	var total int64

	query := s.DB.Model(&entity.ProductCategory{}).Where("organization_id = ?", organizationID)

	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	if err := query.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
