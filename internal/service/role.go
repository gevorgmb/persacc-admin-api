package service

import (
	"context"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type RoleService struct {
	DB *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{DB: db}
}

func (s *RoleService) Create(ctx context.Context, role *entity.Role, permissionIDs []int64) error {
	if len(permissionIDs) > 0 {
		var perms []entity.Permission
		if err := s.DB.Find(&perms, permissionIDs).Error; err != nil {
			return err
		}
		role.Permissions = perms
	}
	return s.DB.Create(role).Error
}

func (s *RoleService) Get(ctx context.Context, id int64) (*entity.Role, error) {
	var role entity.Role
	if err := s.DB.Preload("Permissions").First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) Update(ctx context.Context, role *entity.Role, permissionIDs []int64, updatePerms bool) error {
	if updatePerms {
		var perms []entity.Permission
		if len(permissionIDs) > 0 {
			if err := s.DB.Find(&perms, permissionIDs).Error; err != nil {
				return err
			}
		}

		if err := s.DB.Model(role).Association("Permissions").Replace(perms); err != nil {
			return err
		}
		role.Permissions = perms
	}

	return s.DB.Save(role).Error
}

func (s *RoleService) Delete(ctx context.Context, id int64) error {
	return s.DB.Delete(&entity.Role{}, "id = ?", id).Error
}

func (s *RoleService) List(ctx context.Context, limit, offset int) ([]entity.Role, int64, error) {
	var roles []entity.Role
	var total int64

	s.DB.Model(&entity.Role{}).Count(&total)
	if err := s.DB.Preload("Permissions").Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}
