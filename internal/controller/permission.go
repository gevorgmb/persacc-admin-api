package controller

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/entity"
)

type PermissionController struct {
	DB *gorm.DB
}

func NewPermissionController(db *gorm.DB) *PermissionController {
	return &PermissionController{DB: db}
}

func (c *PermissionController) Create(ctx context.Context, req *adminpb.CreatePermissionRequest) (*adminpb.CreatePermissionResponse, error) {
	permission := entity.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := c.DB.Create(&permission).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create permission: %v", err)
	}

	return &adminpb.CreatePermissionResponse{
		Permission: ConvertPermissionToProto(permission),
	}, nil
}

func (c *PermissionController) Get(ctx context.Context, req *adminpb.GetPermissionRequest) (*adminpb.GetPermissionResponse, error) {
	var permission entity.Permission
	if err := c.DB.First(&permission, "id = ?", req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "permission not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get permission: %v", err)
	}

	return &adminpb.GetPermissionResponse{
		Permission: ConvertPermissionToProto(permission),
	}, nil
}

func (c *PermissionController) Update(ctx context.Context, req *adminpb.UpdatePermissionRequest) (*adminpb.UpdatePermissionResponse, error) {
	var permission entity.Permission
	if err := c.DB.First(&permission, "id = ?", req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "permission not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find permission: %v", err)
	}

	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Description != "" {
		permission.Description = req.Description
	}

	if err := c.DB.Save(&permission).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update permission: %v", err)
	}

	return &adminpb.UpdatePermissionResponse{
		Permission: ConvertPermissionToProto(permission),
	}, nil
}

func (c *PermissionController) Delete(ctx context.Context, req *adminpb.DeletePermissionRequest) (*adminpb.DeletePermissionResponse, error) {
	if err := c.DB.Delete(&entity.Permission{}, "id = ?", req.Id).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete permission: %v", err)
	}
	return &adminpb.DeletePermissionResponse{Success: true}, nil
}

func (c *PermissionController) List(ctx context.Context, req *adminpb.ListPermissionsRequest) (*adminpb.ListPermissionsResponse, error) {
	var permissions []entity.Permission
	var total int64

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	c.DB.Model(&entity.Permission{}).Count(&total)
	if err := c.DB.Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list permissions: %v", err)
	}

	var protoPermissions []*adminpb.Permission
	for _, p := range permissions {
		protoPermissions = append(protoPermissions, ConvertPermissionToProto(p))
	}

	return &adminpb.ListPermissionsResponse{
		Permissions: protoPermissions,
		Total:       int32(total),
		Page:        int32(page),
		Limit:       int32(limit),
	}, nil
}

func ConvertPermissionToProto(p entity.Permission) *adminpb.Permission {
	return &adminpb.Permission{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}
