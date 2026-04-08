package service

import (
	"context"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type SupplierService struct {
	DB *gorm.DB
}

func NewSupplierService(db *gorm.DB) *SupplierService {
	return &SupplierService{DB: db}
}

func (s *SupplierService) Create(ctx context.Context, supplier *entity.Supplier) error {
	return s.DB.Create(supplier).Error
}

func (s *SupplierService) Get(ctx context.Context, id int64, organizationID int64) (*entity.Supplier, error) {
	var supplier entity.Supplier
	err := s.DB.Where("id = ? AND organization_id = ?", id, organizationID).First(&supplier).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (s *SupplierService) Update(ctx context.Context, supplier *entity.Supplier, organizationID int64) error {
	// Verify relationship exists and it matches organization
	var count int64
	s.DB.Model(&entity.Supplier{}).
		Where("id = ? AND organization_id = ?", supplier.ID, organizationID).
		Count(&count)
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return s.DB.Save(supplier).Error
}

func (s *SupplierService) Delete(ctx context.Context, id int64, organizationID int64) error {
	return s.DB.Where("id = ? AND organization_id = ?", id, organizationID).Delete(&entity.Supplier{}).Error
}

func (s *SupplierService) List(ctx context.Context, limit, offset int, organizationID int64, filters map[string]string) ([]entity.Supplier, int64, error) {
	var suppliers []entity.Supplier
	var total int64

	query := s.DB.Model(&entity.Supplier{}).Where("organization_id = ?", organizationID)

	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	if err := query.Limit(limit).Offset(offset).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}
