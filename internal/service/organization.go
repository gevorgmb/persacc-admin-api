package service

import (
	"context"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type OrganizationService struct {
	DB *gorm.DB
}

func NewOrganizationService(db *gorm.DB) *OrganizationService {
	return &OrganizationService{DB: db}
}

func (s *OrganizationService) Create(ctx context.Context, org *entity.Organization) error {
	return s.DB.Create(org).Error
}

func (s *OrganizationService) Get(ctx context.Context, id int64) (*entity.Organization, error) {
	var org entity.Organization
	if err := s.DB.First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (s *OrganizationService) Update(ctx context.Context, org *entity.Organization) error {
	return s.DB.Save(org).Error
}

func (s *OrganizationService) Delete(ctx context.Context, id int64) error {
	return s.DB.Delete(&entity.Organization{}, "id = ?", id).Error
}

func (s *OrganizationService) List(ctx context.Context, limit, offset int) ([]entity.Organization, int64, error) {
	var orgs []entity.Organization
	var total int64

	s.DB.Model(&entity.Organization{}).Count(&total)
	if err := s.DB.Limit(limit).Offset(offset).Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}
