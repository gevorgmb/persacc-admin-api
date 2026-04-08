package service

import (
	"context"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type VendorService struct {
	DB *gorm.DB
}

func NewVendorService(db *gorm.DB) *VendorService {
	return &VendorService{DB: db}
}

func (s *VendorService) Create(ctx context.Context, vendor *entity.Vendor) error {
	return s.DB.Create(vendor).Error
}

func (s *VendorService) Get(ctx context.Context, id int64) (*entity.Vendor, error) {
	var vendor entity.Vendor
	err := s.DB.First(&vendor, id).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (s *VendorService) Update(ctx context.Context, vendor *entity.Vendor) error {
	return s.DB.Save(vendor).Error
}

func (s *VendorService) Delete(ctx context.Context, id int64) error {
	return s.DB.Delete(&entity.Vendor{}, id).Error
}

func (s *VendorService) List(ctx context.Context, limit, offset int, filters map[string]string) ([]entity.Vendor, int64, error) {
	var vendors []entity.Vendor
	var total int64

	query := s.DB.Model(&entity.Vendor{})

	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	if err := query.Limit(limit).Offset(offset).Find(&vendors).Error; err != nil {
		return nil, 0, err
	}

	return vendors, total, nil
}
