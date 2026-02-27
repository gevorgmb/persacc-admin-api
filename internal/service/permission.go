package service

import (
	"context"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type PermissionService struct {
	DB *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{DB: db}
}

func (s *PermissionService) Create(ctx context.Context, permission *entity.Permission) error {
	return s.DB.Create(permission).Error
}

func (s *PermissionService) Get(ctx context.Context, id int64) (*entity.Permission, error) {
	var permission entity.Permission
	if err := s.DB.First(&permission, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (s *PermissionService) Update(ctx context.Context, permission *entity.Permission) error {
	return s.DB.Save(permission).Error
}

func (s *PermissionService) Delete(ctx context.Context, id int64) error {
	return s.DB.Delete(&entity.Permission{}, "id = ?", id).Error
}

func (s *PermissionService) List(ctx context.Context, limit, offset int) ([]entity.Permission, int64, error) {
	var permissions []entity.Permission
	var total int64

	s.DB.Model(&entity.Permission{}).Count(&total)
	if err := s.DB.Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}
