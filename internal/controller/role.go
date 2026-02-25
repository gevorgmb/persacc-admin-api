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

type RoleController struct {
	DB *gorm.DB
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{DB: db}
}

func (c *RoleController) Create(ctx context.Context, req *adminpb.CreateRoleRequest) (*adminpb.CreateRoleResponse, error) {
	role := entity.Role{
		Name: req.Name,
	}

	if len(req.PermissionIds) > 0 {
		var perms []entity.Permission
		if err := c.DB.Find(&perms, req.PermissionIds).Error; err != nil {
			return nil, status.Errorf(codes.Internal, "failed to find permissions: %v", err)
		}
		role.Permissions = perms
	}

	if err := c.DB.Create(&role).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create role: %v", err)
	}

	return &adminpb.CreateRoleResponse{
		Role: ConvertRoleToProto(role),
	}, nil
}

func (c *RoleController) Get(ctx context.Context, req *adminpb.GetRoleRequest) (*adminpb.GetRoleResponse, error) {
	var role entity.Role
	if err := c.DB.Preload("Permissions").First(&role, "id = ?", req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "role not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get role: %v", err)
	}

	return &adminpb.GetRoleResponse{
		Role: ConvertRoleToProto(role),
	}, nil
}

func (c *RoleController) Update(ctx context.Context, req *adminpb.UpdateRoleRequest) (*adminpb.UpdateRoleResponse, error) {
	var role entity.Role
	if err := c.DB.Preload("Permissions").First(&role, "id = ?", req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "role not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find role: %v", err)
	}

	if req.Name != "" {
		role.Name = req.Name
	}

	if req.PermissionIds != nil {
		var perms []entity.Permission
		if len(req.PermissionIds) > 0 {
			if err := c.DB.Find(&perms, req.PermissionIds).Error; err != nil {
				return nil, status.Errorf(codes.Internal, "failed to find permissions: %v", err)
			}
		}

		if err := c.DB.Model(&role).Association("Permissions").Replace(perms); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update role permissions: %v", err)
		}
		role.Permissions = perms
	}

	if err := c.DB.Save(&role).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update role: %v", err)
	}

	return &adminpb.UpdateRoleResponse{
		Role: ConvertRoleToProto(role),
	}, nil
}

func (c *RoleController) Delete(ctx context.Context, req *adminpb.DeleteRoleRequest) (*adminpb.DeleteRoleResponse, error) {
	if err := c.DB.Delete(&entity.Role{}, "id = ?", req.Id).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete role: %v", err)
	}
	return &adminpb.DeleteRoleResponse{Success: true}, nil
}

func (c *RoleController) List(ctx context.Context, req *adminpb.ListRolesRequest) (*adminpb.ListRolesResponse, error) {
	var roles []entity.Role
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

	c.DB.Model(&entity.Role{}).Count(&total)
	if err := c.DB.Preload("Permissions").Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list roles: %v", err)
	}

	var protoRoles []*adminpb.Role
	for _, r := range roles {
		protoRoles = append(protoRoles, ConvertRoleToProto(r))
	}

	return &adminpb.ListRolesResponse{
		Roles: protoRoles,
		Total: int32(total),
		Page:  int32(page),
		Limit: int32(limit),
	}, nil
}

func ConvertRoleToProto(r entity.Role) *adminpb.Role {
	var protoPerms []*adminpb.Permission
	for _, p := range r.Permissions {
		protoPerms = append(protoPerms, ConvertPermissionToProto(p))
	}

	return &adminpb.Role{
		Id:          r.ID,
		Name:        r.Name,
		Permissions: protoPerms,
	}
}
